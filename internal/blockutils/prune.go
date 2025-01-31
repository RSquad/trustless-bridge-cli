package blockutils

import (
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func PruneBlock(blockBOC []byte, isKeyblock bool) (*cell.Cell, error) {
	blockProof, err := BuildBlockProof(blockBOC, isKeyblock)
	if err != nil {
		return nil, err
	}

	return blockProof.PeekRef(0)
}
