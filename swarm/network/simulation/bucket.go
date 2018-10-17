package simulation

import (
	"github.com/TreasureChain/go-tstchain/p2p/discover"
)

// BucketKey is the type that should be used for keys in simulation buckets.
type BucketKey string

// NodeItem returns an item set in ServiceFunc function for a particualar node.
func (s *Simulation) NodeItem(id discover.NodeID, key interface{}) (value interface{}, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.buckets[id]; !ok {
		return nil, false
	}
	return s.buckets[id].Load(key)
}

// SetNodeItem sets a new item associated with the node with provided NodeID.
// Buckets should be used to avoid managing separate simulation global state.
func (s *Simulation) SetNodeItem(id discover.NodeID, key interface{}, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.buckets[id].Store(key, value)
}

// NodesItems returns a map of items from all nodes that are all set under the
// same BucketKey.
func (s *Simulation) NodesItems(key interface{}) (values map[discover.NodeID]interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.NodeIDs()
	values = make(map[discover.NodeID]interface{}, len(ids))
	for _, id := range ids {
		if _, ok := s.buckets[id]; !ok {
			continue
		}
		if v, ok := s.buckets[id].Load(key); ok {
			values[id] = v
		}
	}
	return values
}

// UpNodesItems returns a map of items with the same BucketKey from all nodes that are up.
func (s *Simulation) UpNodesItems(key interface{}) (values map[discover.NodeID]interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.UpNodeIDs()
	values = make(map[discover.NodeID]interface{})
	for _, id := range ids {
		if _, ok := s.buckets[id]; !ok {
			continue
		}
		if v, ok := s.buckets[id].Load(key); ok {
			values[id] = v
		}
	}
	return values
}
