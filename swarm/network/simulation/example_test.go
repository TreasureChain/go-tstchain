package simulation_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/TreasureChain/go-tstchain/log"
	"github.com/TreasureChain/go-tstchain/node"
	"github.com/TreasureChain/go-tstchain/p2p"
	"github.com/TreasureChain/go-tstchain/p2p/simulations/adapters"
	"github.com/TreasureChain/go-tstchain/swarm/network"
	"github.com/TreasureChain/go-tstchain/swarm/network/simulation"
)

// Every node can have a Kademlia associated using the node bucket under
// BucketKeyKademlia key. This allows to use WaitTillHealthy to block until
// all nodes have the their Kadmlias healthy.
func ExampleSimulation_WaitTillHealthy() {
	sim := simulation.New(map[string]simulation.ServiceFunc{
		"bzz": func(ctx *adapters.ServiceContext, b *sync.Map) (node.Service, func(), error) {
			addr := network.NewAddrFromNodeID(ctx.Config.ID)
			hp := network.NewHiveParams()
			hp.Discovery = false
			config := &network.BzzConfig{
				OverlayAddr:  addr.Over(),
				UnderlayAddr: addr.Under(),
				HiveParams:   hp,
			}
			kad := network.NewKademlia(addr.Over(), network.NewKadParams())
			// store kademlia in node's bucket under BucketKeyKademlia
			// so that it can be found by WaitTillHealthy method.
			b.Store(simulation.BucketKeyKademlia, kad)
			return network.NewBzz(config, kad, nil, nil, nil), nil, nil
		},
	})
	defer sim.Close()

	_, err := sim.AddNodesAndConnectRing(10)
	if err != nil {
		// handle error properly...
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	ill, err := sim.WaitTillHealthy(ctx, 2)
	if err != nil {
		// inspect the latest detected not healthy kademlias
		for id, kad := range ill {
			fmt.Println("Node", id)
			fmt.Println(kad.String())
		}
		// handle error...
	}

	// continue with the test
}

// Watch all peer events in the simulation network, buy receiving from a channel.
func ExampleSimulation_PeerEvents() {
	sim := simulation.New(nil)
	defer sim.Close()

	events := sim.PeerEvents(context.Background(), sim.NodeIDs())

	go func() {
		for e := range events {
			if e.Error != nil {
				log.Error("peer event", "err", e.Error)
				continue
			}
			log.Info("peer event", "node", e.NodeID, "peer", e.Event.Peer, "msgcode", e.Event.MsgCode)
		}
	}()
}

// Detect when a nodes drop a peer.
func ExampleSimulation_PeerEvents_disconnections() {
	sim := simulation.New(nil)
	defer sim.Close()

	disconnections := sim.PeerEvents(
		context.Background(),
		sim.NodeIDs(),
		simulation.NewPeerEventsFilter().Type(p2p.PeerEventTypeDrop),
	)

	go func() {
		for d := range disconnections {
			if d.Error != nil {
				log.Error("peer drop", "err", d.Error)
				continue
			}
			log.Warn("peer drop", "node", d.NodeID, "peer", d.Event.Peer)
		}
	}()
}

// Watch multiple types of events or messages. In this case, they differ only
// by MsgCode, but filters can be set for different types or protocols, too.
func ExampleSimulation_PeerEvents_multipleFilters() {
	sim := simulation.New(nil)
	defer sim.Close()

	msgs := sim.PeerEvents(
		context.Background(),
		sim.NodeIDs(),
		// Watch when bzz messages 1 and 4 are received.
		simulation.NewPeerEventsFilter().Type(p2p.PeerEventTypeMsgRecv).Protocol("bzz").MsgCode(1),
		simulation.NewPeerEventsFilter().Type(p2p.PeerEventTypeMsgRecv).Protocol("bzz").MsgCode(4),
	)

	go func() {
		for m := range msgs {
			if m.Error != nil {
				log.Error("bzz message", "err", m.Error)
				continue
			}
			log.Info("bzz message", "node", m.NodeID, "peer", m.Event.Peer)
		}
	}()
}
