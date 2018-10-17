package mem

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"

	"github.com/TreasureChain/go-tstchain/common"
	"github.com/TreasureChain/go-tstchain/swarm/storage/mock"
)

// GlobalStore stores all chunk data and also keys and node addresses relations.
// It implements mock.GlobalStore interface.
type GlobalStore struct {
	nodes map[string]map[common.Address]struct{}
	data  map[string][]byte
	mu    sync.Mutex
}

// NewGlobalStore creates a new instance of GlobalStore.
func NewGlobalStore() *GlobalStore {
	return &GlobalStore{
		nodes: make(map[string]map[common.Address]struct{}),
		data:  make(map[string][]byte),
	}
}

// NewNodeStore returns a new instance of NodeStore that retrieves and stores
// chunk data only for a node with address addr.
func (s *GlobalStore) NewNodeStore(addr common.Address) *mock.NodeStore {
	return mock.NewNodeStore(addr, s)
}

// Get returns chunk data if the chunk with key exists for node
// on address addr.
func (s *GlobalStore) Get(addr common.Address, key []byte) (data []byte, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.nodes[string(key)][addr]; !ok {
		return nil, mock.ErrNotFound
	}

	data, ok := s.data[string(key)]
	if !ok {
		return nil, mock.ErrNotFound
	}
	return data, nil
}

// Put saves the chunk data for node with address addr.
func (s *GlobalStore) Put(addr common.Address, key []byte, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.nodes[string(key)]; !ok {
		s.nodes[string(key)] = make(map[common.Address]struct{})
	}
	s.nodes[string(key)][addr] = struct{}{}
	s.data[string(key)] = data
	return nil
}

// HasKey returns whether a node with addr contains the key.
func (s *GlobalStore) HasKey(addr common.Address, key []byte) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.nodes[string(key)][addr]
	return ok
}

// Import reads tar archive from a reader that contains exported chunk data.
// It returns the number of chunks imported and an error.
func (s *GlobalStore) Import(r io.Reader) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tr := tar.NewReader(r)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return n, err
		}

		data, err := ioutil.ReadAll(tr)
		if err != nil {
			return n, err
		}

		var c mock.ExportedChunk
		if err = json.Unmarshal(data, &c); err != nil {
			return n, err
		}

		addrs := make(map[common.Address]struct{})
		for _, a := range c.Addrs {
			addrs[a] = struct{}{}
		}

		key := string(common.Hex2Bytes(hdr.Name))
		s.nodes[key] = addrs
		s.data[key] = c.Data
		n++
	}
	return n, err
}

// Export writes to a writer a tar archive with all chunk data from
// the store. It returns the number of chunks exported and an error.
func (s *GlobalStore) Export(w io.Writer) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tw := tar.NewWriter(w)
	defer tw.Close()

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := json.NewEncoder(buf)
	for key, addrs := range s.nodes {
		al := make([]common.Address, 0, len(addrs))
		for a := range addrs {
			al = append(al, a)
		}

		buf.Reset()
		if err = encoder.Encode(mock.ExportedChunk{
			Addrs: al,
			Data:  s.data[key],
		}); err != nil {
			return n, err
		}

		data := buf.Bytes()
		hdr := &tar.Header{
			Name: common.Bytes2Hex([]byte(key)),
			Mode: 0644,
			Size: int64(len(data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return n, err
		}
		if _, err := tw.Write(data); err != nil {
			return n, err
		}
		n++
	}
	return n, err
}