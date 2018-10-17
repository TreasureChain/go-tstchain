package mock

import (
	"errors"
	"io"

	"github.com/TreasureChain/go-tstchain/common"
)

// ErrNotFound indicates that the chunk is not found.
var ErrNotFound = errors.New("not found")

// NodeStore holds the node address and a reference to the GlobalStore
// in order to access and store chunk data only for one node.
type NodeStore struct {
	store GlobalStorer
	addr  common.Address
}

// NewNodeStore creates a new instance of NodeStore that keeps
// chunk data using GlobalStorer with a provided address.
func NewNodeStore(addr common.Address, store GlobalStorer) *NodeStore {
	return &NodeStore{
		store: store,
		addr:  addr,
	}
}

// Get returns chunk data for a key for a node that has the address
// provided on NodeStore initialization.
func (n *NodeStore) Get(key []byte) (data []byte, err error) {
	return n.store.Get(n.addr, key)
}

// Put saves chunk data for a key for a node that has the address
// provided on NodeStore initialization.
func (n *NodeStore) Put(key []byte, data []byte) error {
	return n.store.Put(n.addr, key, data)
}

// GlobalStorer defines methods for mock db store
// that stores chunk data for all swarm nodes.
// It is used in tests to construct mock NodeStores
// for swarm nodes and to track and validate chunks.
type GlobalStorer interface {
	Get(addr common.Address, key []byte) (data []byte, err error)
	Put(addr common.Address, key []byte, data []byte) error
	HasKey(addr common.Address, key []byte) bool
	// NewNodeStore creates an instance of NodeStore
	// to be used by a single swarm node with
	// address addr.
	NewNodeStore(addr common.Address) *NodeStore
}

// Importer defines method for importing mock store data
// from an exported tar archive.
type Importer interface {
	Import(r io.Reader) (n int, err error)
}

// Exporter defines method for exporting mock store data
// to a tar archive.
type Exporter interface {
	Export(w io.Writer) (n int, err error)
}

// ImportExporter is an interface for importing and exporting
// mock store data to and from a tar archive.
type ImportExporter interface {
	Importer
	Exporter
}

// ExportedChunk is the structure that is saved in tar archive for
// each chunk as JSON-encoded bytes.
type ExportedChunk struct {
	Data  []byte           `json:"d"`
	Addrs []common.Address `json:"a"`
}
