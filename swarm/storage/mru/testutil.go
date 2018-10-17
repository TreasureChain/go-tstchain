package mru

import (
	"fmt"
	"path/filepath"

	"github.com/TreasureChain/go-tstchain/swarm/storage"
)

const (
	testDbDirName = "mru"
)

type TestHandler struct {
	*Handler
}

func (t *TestHandler) Close() {
	t.chunkStore.Close()
}

// NewTestHandler creates Handler object to be used for testing purposes.
func NewTestHandler(datadir string, params *HandlerParams) (*TestHandler, error) {
	path := filepath.Join(datadir, testDbDirName)
	rh, err := NewHandler(params)
	if err != nil {
		return nil, fmt.Errorf("resource handler create fail: %v", err)
	}
	localstoreparams := storage.NewDefaultLocalStoreParams()
	localstoreparams.Init(path)
	localStore, err := storage.NewLocalStore(localstoreparams, nil)
	if err != nil {
		return nil, fmt.Errorf("localstore create fail, path %s: %v", path, err)
	}
	localStore.Validators = append(localStore.Validators, storage.NewContentAddressValidator(storage.MakeHashFunc(resourceHashAlgorithm)))
	localStore.Validators = append(localStore.Validators, rh)
	netStore := storage.NewNetStore(localStore, nil)
	rh.SetStore(netStore)
	return &TestHandler{rh}, nil
}
