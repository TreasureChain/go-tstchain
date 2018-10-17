package ethclient

import "github.com/TreasureChain/go-tstchain"

// Verify that Client implements the tstchain interfaces.
var (
	_ = tstchain.ChainReader(&Client{})
	_ = tstchain.TransactionReader(&Client{})
	_ = tstchain.ChainStateReader(&Client{})
	_ = tstchain.ChainSyncReader(&Client{})
	_ = tstchain.ContractCaller(&Client{})
	_ = tstchain.GasEstimator(&Client{})
	_ = tstchain.GasPricer(&Client{})
	_ = tstchain.LogFilterer(&Client{})
	_ = tstchain.PendingStateReader(&Client{})
	// _ = tstchain.PendingStateEventer(&Client{})
	_ = tstchain.PendingContractCaller(&Client{})
)
