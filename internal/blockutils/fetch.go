package blockutils

import (
	"context"
	"fmt"

	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

func FetchMasterchainBlock(ctx context.Context, tonClient *tonclient.TonClient, seqno uint32) (*tlb.Block, error) {
	blockIDExt, err := tonClient.API.LookupBlock(ctx, -1, 0, seqno)
	if err != nil {
		return nil, err
	}

	return tonClient.API.GetBlockData(ctx, blockIDExt)
}

func FetchMasterchainBlockBOC(
	ctx context.Context,
	tonClient *tonclient.TonClient,
	seqno uint32,
) (*ton.BlockIDExt, []byte, error) {
	blockIDExt, err := tonClient.API.LookupBlock(ctx, -1, 0, seqno)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup block: %w", err)
	}
	blockBOC, err := tonClient.GetBlockBOC(ctx, blockIDExt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get block BOC: %w", err)
	}
	return blockIDExt, blockBOC, nil
}
