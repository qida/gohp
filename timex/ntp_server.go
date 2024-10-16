package timex

import (
	"github.com/btfak/sntp/netapp"
	"github.com/btfak/sntp/netevent"
)

// 默认123
func ServerNtp(port int) {
	handler := netapp.GetHandler()
	netevent.Reactor.ListenUdp(port, handler)
	netevent.Reactor.Run()
}
