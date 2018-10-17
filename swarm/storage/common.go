package storage

import (
	"context"
	"sync"

	"github.com/TreasureChain/go-tstchain/swarm/log"
)

// PutChunks adds chunks  to localstore
// It waits for receive on the stored channel
// It logs but does not fail on delivery error
func PutChunks(store *LocalStore, chunks ...*Chunk) {
	wg := sync.WaitGroup{}
	wg.Add(len(chunks))
	go func() {
		for _, c := range chunks {
			<-c.dbStoredC
			if err := c.GetErrored(); err != nil {
				log.Error("chunk store fail", "err", err, "key", c.Addr)
			}
			wg.Done()
		}
	}()
	for _, c := range chunks {
		go store.Put(context.TODO(), c)
	}
	wg.Wait()
}
