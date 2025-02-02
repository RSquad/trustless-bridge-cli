package blockutils

import (
	"context"

	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/xssnick/tonutils-go/tlb"
)

func FetchMasterchainBlock(ctx context.Context, tonClient *tonclient.TonClient, seqno uint32) (*tlb.Block, error) {
	blockIDExt, err := tonClient.API.LookupBlock(context.Background(), -1, 0, seqno)
	if err != nil {
		return nil, err
	}

	return tonClient.API.GetBlockData(context.Background(), blockIDExt)
}
