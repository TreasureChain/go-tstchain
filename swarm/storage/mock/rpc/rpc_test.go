package rpc

import (
	"testing"

	"github.com/TreasureChain/go-tstchain/rpc"
	"github.com/TreasureChain/go-tstchain/swarm/storage/mock/mem"
	"github.com/TreasureChain/go-tstchain/swarm/storage/mock/test"
)

// TestDBStore is running test for a GlobalStore
// using test.MockStore function.
func TestRPCStore(t *testing.T) {
	serverStore := mem.NewGlobalStore()

	server := rpc.NewServer()
	if err := server.RegisterName("mockStore", serverStore); err != nil {
		t.Fatal(err)
	}

	store := NewGlobalStore(rpc.DialInProc(server))
	defer store.Close()

	test.MockStore(t, store, 100)
}
