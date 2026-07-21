package hysteria

import (
	"github.com/sagernet/sing-box/option"
)

func (h *Inbound) UpdateUsers(users []option.HysteriaUser) error {
	userList := make([]int, 0, len(users))
	userNameList := make([]string, 0, len(users))
	userPasswordList := make([]string, 0, len(users))
	for index, user := range users {
		userList = append(userList, index)
		userNameList = append(userNameList, user.Name)
		var password string
		if user.AuthString != "" {
			password = user.AuthString
		} else {
			password = string(user.Auth)
		}
		userPasswordList = append(userPasswordList, password)
	}
	// Grow before the service learns the new users, shrink after it forgets the
	// old ones. The service hands back an index into userNameList on every new
	// connection, so it must never know about more users than the slice holds
	// or that index goes out of range. (The two are still unsynchronised, as
	// they are in sing-box itself -- a racing read can see a stale name, which
	// only mis-labels a stat. Out of range would take down the process.)
	if len(userNameList) > len(h.userNameList) {
		h.userNameList = userNameList
	}
	h.service.UpdateUsers(userList, userPasswordList)
	h.userNameList = userNameList
	return nil
}
