package vmess

import (
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func (h *Inbound) UpdateUsers(users []option.VMessUser) error {
	// Grow before the service learns the new users, shrink after. See the
	// comment in core/protocol/hysteria/users.go for why the order matters.
	if len(users) > len(h.users) {
		h.users = users
	}
	err := h.service.UpdateUsers(common.MapIndexed(users, func(index int, _ option.VMessUser) int {
		return index
	}), common.Map(users, func(it option.VMessUser) string {
		return it.UUID
	}), common.Map(users, func(it option.VMessUser) int {
		return it.AlterId
	}))
	if err != nil {
		return err
	}
	h.users = users
	return nil
}
