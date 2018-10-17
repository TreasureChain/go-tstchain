package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TreasureChain/go-tstchain/common/hexutil"

	"github.com/TreasureChain/go-tstchain/cmd/utils"
	swarm "github.com/TreasureChain/go-tstchain/swarm/api/client"
	"github.com/TreasureChain/go-tstchain/swarm/storage/mru"
	"gopkg.in/urfave/cli.v1"
)

func NewGenericSigner(ctx *cli.Context) mru.Signer {
	return mru.NewGenericSigner(getPrivKey(ctx))
}

// swarm resource create <frequency> [--name <name>] [--data <0x Hexdata> [--multihash=false]]
// swarm resource update <Manifest Address or ENS domain> <0x Hexdata> [--multihash=false]
// swarm resource info <Manifest Address or ENS domain>

func resourceCreate(ctx *cli.Context) {
	args := ctx.Args()

	var (
		bzzapi      = strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
		client      = swarm.NewClient(bzzapi)
		multihash   = ctx.Bool(SwarmResourceMultihashFlag.Name)
		initialData = ctx.String(SwarmResourceDataOnCreateFlag.Name)
		name        = ctx.String(SwarmResourceNameFlag.Name)
	)

	if len(args) < 1 {
		fmt.Println("Incorrect number of arguments")
		cli.ShowCommandHelpAndExit(ctx, "create", 1)
		return
	}
	signer := NewGenericSigner(ctx)
	frequency, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Printf("Frequency formatting error: %s\n", err.Error())
		cli.ShowCommandHelpAndExit(ctx, "create", 1)
		return
	}

	metadata := mru.ResourceMetadata{
		Name:      name,
		Frequency: frequency,
		Owner:     signer.Address(),
	}

	var newResourceRequest *mru.Request
	if initialData != "" {
		initialDataBytes, err := hexutil.Decode(initialData)
		if err != nil {
			fmt.Printf("Error parsing data: %s\n", err.Error())
			cli.ShowCommandHelpAndExit(ctx, "create", 1)
			return
		}
		newResourceRequest, err = mru.NewCreateUpdateRequest(&metadata)
		if err != nil {
			utils.Fatalf("Error creating new resource request: %s", err)
		}
		newResourceRequest.SetData(initialDataBytes, multihash)
		if err = newResourceRequest.Sign(signer); err != nil {
			utils.Fatalf("Error signing resource update: %s", err.Error())
		}
	} else {
		newResourceRequest, err = mru.NewCreateRequest(&metadata)
		if err != nil {
			utils.Fatalf("Error creating new resource request: %s", err)
		}
	}

	manifestAddress, err := client.CreateResource(newResourceRequest)
	if err != nil {
		utils.Fatalf("Error creating resource: %s", err.Error())
		return
	}
	fmt.Println(manifestAddress) // output manifest address to the user in a single line (useful for other commands to pick up)

}

func resourceUpdate(ctx *cli.Context) {
	args := ctx.Args()

	var (
		bzzapi    = strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
		client    = swarm.NewClient(bzzapi)
		multihash = ctx.Bool(SwarmResourceMultihashFlag.Name)
	)

	if len(args) < 2 {
		fmt.Println("Incorrect number of arguments")
		cli.ShowCommandHelpAndExit(ctx, "update", 1)
		return
	}
	signer := NewGenericSigner(ctx)
	manifestAddressOrDomain := args[0]
	data, err := hexutil.Decode(args[1])
	if err != nil {
		utils.Fatalf("Error parsing data: %s", err.Error())
		return
	}

	// Retrieve resource status and metadata out of the manifest
	updateRequest, err := client.GetResourceMetadata(manifestAddressOrDomain)
	if err != nil {
		utils.Fatalf("Error retrieving resource status: %s", err.Error())
	}

	// set the new data
	updateRequest.SetData(data, multihash)

	// sign update
	if err = updateRequest.Sign(signer); err != nil {
		utils.Fatalf("Error signing resource update: %s", err.Error())
	}

	// post update
	err = client.UpdateResource(updateRequest)
	if err != nil {
		utils.Fatalf("Error updating resource: %s", err.Error())
		return
	}
}

func resourceInfo(ctx *cli.Context) {
	var (
		bzzapi = strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
		client = swarm.NewClient(bzzapi)
	)
	args := ctx.Args()
	if len(args) < 1 {
		fmt.Println("Incorrect number of arguments.")
		cli.ShowCommandHelpAndExit(ctx, "info", 1)
		return
	}
	manifestAddressOrDomain := args[0]
	metadata, err := client.GetResourceMetadata(manifestAddressOrDomain)
	if err != nil {
		utils.Fatalf("Error retrieving resource metadata: %s", err.Error())
		return
	}
	encodedMetadata, err := metadata.MarshalJSON()
	if err != nil {
		utils.Fatalf("Error encoding metadata to JSON for display:%s", err)
	}
	fmt.Println(string(encodedMetadata))
}
