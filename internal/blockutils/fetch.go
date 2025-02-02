package blockutils

import (
	"context"
	"fmt"

	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func FetchMasterchainBlock(ctx context.Context, tonClient *tonclient.TonClient, seqno uint32) (*tlb.Block, error) {
	blockIDExt, err := tonClient.API.LookupBlock(ctx, -1, 0, seqno)
	if err != nil {
		return nil, err
	}

	return tonClient.API.GetBlockData(ctx, blockIDExt)
}

func FetchMasterchainBlockCell(
	ctx context.Context,
	tonClient *tonclient.TonClient,
	seqno uint32,
) (*cell.Cell, error) {
	blockIDExt, err := tonClient.API.LookupBlock(ctx, -1, 0, seqno)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup block: %w", err)
	}
	blockBOC, err := tonClient.GetBlockBOC(ctx, blockIDExt)
	if err != nil {
		return nil, fmt.Errorf("failed to get block BOC: %w", err)
	}
	blockCell, err := cell.FromBOC(blockBOC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse block BOC: %w", err)
	}
	return blockCell, nil
}
