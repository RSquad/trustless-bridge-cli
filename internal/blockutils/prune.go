package blockutils

import (
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func PruneBlock(blockBOC []byte) (*cell.Cell, error) {
	blockProof, err := BuildBlockProof(blockBOC)
	if err != nil {
		return nil, err
	}

	return blockProof.PeekRef(0)
}
