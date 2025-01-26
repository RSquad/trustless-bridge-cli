package blockutils

import (
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func BuildBlockProof(blockBOC []byte) (*cell.Cell, error) {
	blockCell, err := cell.FromBOC(blockBOC)
	if err != nil {
		return nil, err
	}

	isCustomCellExists := true
	_, err = blockCell.MustPeekRef(3).PeekRef(3)
	if err != nil {
		isCustomCellExists = false
	}

	sk := createBlockProofSk(isCustomCellExists)

	blockProof, err := blockCell.CreateProof(sk)
	if err != nil {
		return nil, err
	}

	return blockProof, nil
}

func createBlockProofSk(isCustomCellExists bool) *cell.ProofSkeleton {
	sk := cell.CreateProofSkeleton()
	extraSk := sk.ProofRef(3)
	extraSk.ProofRef(2).SetRecursive() // account_blocks
	if isCustomCellExists {
		extraSk.ProofRef(3).SetRecursive() // custom
	}
	return sk
}
