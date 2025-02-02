package tonclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/rsquad/trustless-bridge-cli/internal/data"
	"github.com/spf13/viper"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tl"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type TonClient struct {
	connPool *liteclient.ConnectionPool
	API      *ton.APIClient
}

func NewTonClient(cfg *liteclient.GlobalConfig) (*TonClient, error) {
	connPool := liteclient.NewConnectionPool()

	err := connPool.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	apiWrapped := ton.NewAPIClient(connPool).WithRetry(3)
	api, _ := apiWrapped.(*ton.APIClient)

	return &TonClient{connPool: connPool, API: api}, nil
}

func NewTonClientNetwork(network string) (*TonClient, error) {
	var configDataStr string
	switch network {
	case "testnet":
		configDataStr = data.TestnetConfig
	case "fastnet":
		configDataStr = data.FastnetConfig
	default:
		log.Fatalf("unknown network: %s", network)
	}

	var globalConfig liteclient.GlobalConfig
	err := json.Unmarshal([]byte(configDataStr), &globalConfig)
	if err != nil {
		log.Fatalf("failed to parse config data: %v", err)
	}

	return NewTonClient(&globalConfig)
}

func (tc *TonClient) GetBlockProofExt(ctx context.Context, known, target *ton.BlockIDExt) (*ton.PartialBlockProof, error) {
	var resp tl.Serializable
	err := tc.API.Client().QueryLiteserver(ctx, ton.GetBlockProof{
		Mode:        0x1001,
		KnownBlock:  known,
		TargetBlock: target,
	}, &resp)
	if err != nil {
		return nil, err
	}

	switch t := resp.(type) {
	case ton.PartialBlockProof:
		return &t, nil
	case ton.LSError:
		return nil, t
	}
	return nil, fmt.Errorf("unknown response type")
}

func (tc *TonClient) GetBlockBOC(ctx context.Context, block *ton.BlockIDExt) ([]byte, error) {
	var resp tl.Serializable
	err := tc.API.Client().QueryLiteserver(ctx, ton.GetBlockData{ID: block}, &resp)
	if err != nil {
		return nil, err
	}

	switch t := resp.(type) {
	case ton.BlockData:
		return t.Payload, nil
	case ton.LSError:
		return nil, t
	}
	panic("should not happen")
}

func (tc *TonClient) GetWallet() *wallet.Wallet {
	mnemonic := viper.GetString("wallet_mnemonic")
	if mnemonic == "" {
		panic("wallet_mnemonic is not set")
	}
	walletVersion := viper.GetString("wallet_version")
	if walletVersion == "" {
		panic("wallet_version is not set")
	}

	versionMap := map[string]wallet.Version{
		"v1r1":      wallet.V1R1,
		"v1r2":      wallet.V1R2,
		"v1r3":      wallet.V1R3,
		"v2r1":      wallet.V2R1,
		"v2r2":      wallet.V2R2,
		"v3r1":      wallet.V3R1,
		"v3r2":      wallet.V3R2,
		"v3":        wallet.V3,
		"v4r1":      wallet.V4R1,
		"v4r2":      wallet.V4R2,
		"v5r1beta":  wallet.V5R1Beta,
		"v5r1final": wallet.V5R1Final,
	}

	version, exists := versionMap[strings.ToLower(walletVersion)]
	if !exists {
		panic(fmt.Sprintf("unsupported wallet type: %s", walletVersion))
	}

	w, err := wallet.FromSeed(tc.API, strings.Split(mnemonic, " "), version)
	if err != nil {
		panic(err)
	}
	return w
}
