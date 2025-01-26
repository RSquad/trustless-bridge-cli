package tonclient

import (
	"context"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tl"
	"github.com/xssnick/tonutils-go/ton"
)

type TonClient struct {
	connPool *liteclient.ConnectionPool
	API      *ton.APIClient
}

func NewTonClient(configUrl string) (*TonClient, error) {
	connPool := liteclient.NewConnectionPool()

	err := connPool.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return nil, err
	}
	apiWrapped := ton.NewAPIClient(connPool).WithRetry(3)
	api, _ := apiWrapped.(*ton.APIClient)

	return &TonClient{connPool: connPool, API: api}, nil
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
