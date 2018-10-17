package simulation

import (
	"github.com/TreasureChain/go-tstchain/node"
	"github.com/TreasureChain/go-tstchain/p2p/discover"
	"github.com/TreasureChain/go-tstchain/p2p/simulations/adapters"
)

// Service returns a single Service by name on a particular node
// with provided id.
func (s *Simulation) Service(name string, id discover.NodeID) node.Service {
	simNode, ok := s.Net.GetNode(id).Node.(*adapters.SimNode)
	if !ok {
		return nil
	}
	services := simNode.ServiceMap()
	if len(services) == 0 {
		return nil
	}
	return services[name]
}

// RandomService returns a single Service by name on a
// randomly chosen node that is up.
func (s *Simulation) RandomService(name string) node.Service {
	n := s.RandomUpNode()
	if n == nil {
		return nil
	}
	return n.Service(name)
}

// Services returns all services with a provided name
// from nodes that are up.
func (s *Simulation) Services(name string) (services map[discover.NodeID]node.Service) {
	nodes := s.Net.GetNodes()
	services = make(map[discover.NodeID]node.Service)
	for _, node := range nodes {
		if !node.Up {
			continue
		}
		simNode, ok := node.Node.(*adapters.SimNode)
		if !ok {
			continue
		}
		services[node.ID()] = simNode.Service(name)
	}
	return services
}
