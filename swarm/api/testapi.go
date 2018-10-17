package api

import (
	"github.com/TreasureChain/go-tstchain/swarm/network"
)

type Control struct {
	api  *API
	hive *network.Hive
}

func NewControl(api *API, hive *network.Hive) *Control {
	return &Control{api, hive}
}

//func (self *Control) BlockNetworkRead(on bool) {
//	self.hive.BlockNetworkRead(on)
//}
//
//func (self *Control) SyncEnabled(on bool) {
//	self.hive.SyncEnabled(on)
//}
//
//func (self *Control) SwapEnabled(on bool) {
//	self.hive.SwapEnabled(on)
//}
//
func (c *Control) Hive() string {
	return c.hive.String()
}
