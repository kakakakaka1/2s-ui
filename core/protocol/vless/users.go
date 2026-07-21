package vless

import (
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func (h *Inbound) UpdateUsers(users []option.VLESSUser) error {
	// Grow before the service learns the new users, shrink after. See the
	// comment in core/protocol/hysteria/users.go for why the order matters.
	if len(users) > len(h.users) {
		h.users = users
	}
	h.service.UpdateUsers(common.MapIndexed(users, func(index int, _ option.VLESSUser) int {
		return index
	}), common.Map(users, func(it option.VLESSUser) string {
		return it.UUID
	}), common.Map(users, func(it option.VLESSUser) string {
		return it.Flow
	}))
	h.users = users
	return nil
}
