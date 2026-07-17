package migration

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shenaba/2s-ui/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// versionBefore reports whether dot-separated numeric version a sorts before
// b. The gates below used plain string compares, which break once a segment
// reaches two digits ("1.5.10" < "1.5.7" lexically). Missing or non-numeric
// segments count as 0, so "1.5" == "1.5.0" and "" sorts first.
func versionBefore(a, b string) bool {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		av, bv := 0, 0
		if i < len(as) {
			av, _ = strconv.Atoi(as[i])
		}
		if i < len(bs) {
			bv, _ = strconv.Atoi(bs[i])
		}
		if av != bv {
			return av < bv
		}
	}
	return false
}

func MigrateDb() {
	// void running on first install
	path := config.GetDBPath()
	_, err := os.Stat(path)
	if err != nil {
		println("Database not found")
		return
	}

	db, err := gorm.Open(sqlite.Open(path))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if sqlDB, e := db.DB(); e == nil {
			_ = sqlDB.Close()
		}
	}()
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	currentVersion := config.GetVersion()
	dbVersion := ""
	tx.Raw("SELECT value FROM settings WHERE key = ?", "version").Find(&dbVersion)
	fmt.Println("Current version:", currentVersion, "\nDatabase version:", dbVersion)

	if currentVersion == dbVersion {
		fmt.Println("Database is up to date, no need to migrate")
		return
	}

	fmt.Println("Start migrating database...")

	// Before 1.2
	if dbVersion == "" {
		err = to1_1(tx)
		if err != nil {
			log.Fatal("Migration to 1.1 failed: ", err)
			return
		}
		err = to1_2(tx)
		if err != nil {
			log.Fatal("Migration to 1.2 failed: ", err)
			return
		}
		dbVersion = "1.2"
	}

	// Before 1.3
	if dbVersion[0:3] == "1.2" {
		err = to1_3(tx)
		if err != nil {
			log.Fatal("Migration to 1.3 failed: ", err)
			return
		}
	}

	// 2s-ui version line: both upstream migrations below first ship with
	// 2s-ui 1.5.4, and our users' dbVersion follows 2s-ui releases (1.4.2 ..
	// 1.5.3), NOT upstream's. Gate on 1.5.4, not the upstream version numbers,
	// or every existing install would skip them. Both are idempotent.

	// Back-fill self-signed TLS public-key pins and rewrite OutJson
	if versionBefore(dbVersion, "1.5.4") {
		err = to1_5_1(tx)
		if err != nil {
			log.Fatal("Migration to 1.5.1 failed: ", err)
			return
		}
	}

	// Hash any plaintext admin passwords
	if versionBefore(dbVersion, "1.5.4") {
		err = to1_5_2(tx)
		if err != nil {
			log.Fatal("Migration to 1.5.2 failed: ", err)
			return
		}
	}

	// Strip server-only TLS fields leaked into client-facing JSON (#51)
	if versionBefore(dbVersion, "1.5.7") {
		err = to1_5_7(tx)
		if err != nil {
			log.Fatal("Migration to 1.5.7 failed: ", err)
			return
		}
	}

	// Set version
	err = tx.Exec("UPDATE settings SET value = ? WHERE key = ?", currentVersion, "version").Error
	if err != nil {
		log.Fatal("Update version failed: ", err)
		return
	}
	fmt.Println("Migration done!")
}
