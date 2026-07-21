package hysteria2

import (
	"github.com/sagernet/sing-box/option"
)

func (h *Inbound) UpdateUsers(users []option.Hysteria2User) error {
	userList := make([]int, 0, len(users))
	userNameList := make([]string, 0, len(users))
	userPasswordList := make([]string, 0, len(users))
	for index, user := range users {
		userList = append(userList, index)
		userNameList = append(userNameList, user.Name)
		userPasswordList = append(userPasswordList, user.Password)
	}
	// Grow before the service learns the new users, shrink after. See the
	// comment in core/protocol/hysteria/users.go for why the order matters.
	if len(userNameList) > len(h.userNameList) {
		h.userNameList = userNameList
	}
	h.service.UpdateUsers(userList, userPasswordList)
	h.userNameList = userNameList
	return nil
}
