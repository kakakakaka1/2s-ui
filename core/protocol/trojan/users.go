package trojan

import (
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func (h *Inbound) UpdateUsers(users []option.TrojanUser) error {
	// Grow before the service learns the new users, shrink after. See the
	// comment in core/protocol/hysteria/users.go for why the order matters.
	// On a failed update the list stays grown, which is harmless: every index
	// the service can still produce is in range, only the name may be stale.
	if len(users) > len(h.users) {
		h.users = users
	}
	err := h.service.UpdateUsers(common.MapIndexed(users, func(index int, _ option.TrojanUser) int {
		return index
	}), common.Map(users, func(it option.TrojanUser) string {
		return it.Password
	}))
	if err != nil {
		return err
	}
	h.users = users
	return nil
}
