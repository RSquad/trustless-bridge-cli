package blockutils

import (
	"fmt"

	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func BuildBlockProof(blockBOC []byte, isKeyblock bool) (*cell.Cell, error) {
	blockCell, err := cell.FromBOC(blockBOC)
	if err != nil {
		return nil, err
	}

	var block tlb.Block
	err = tlb.LoadFromCell(&block, blockCell.BeginParse())
	if err != nil {
		return nil, err
	}

	if isKeyblock {
		if block.Extra == nil || block.Extra.Custom == nil || block.Extra.Custom.ConfigParams == nil {
			return nil, fmt.Errorf("extra or custom or config params is nil")
		}
		root, sk := createKeyBlockProofSk()
		_, configSk, err := block.Extra.Custom.ConfigParams.Config.Params.LoadValueWithProof(
			cell.BeginCell().MustStoreUInt(34, 32).EndCell(),
			sk,
		)
		if err != nil {
			return nil, err
		}
		configSk.SetRecursive()

		return blockCell.CreateProof(root)
	}

	sk := createBlockProofSk()
	return blockCell.CreateProof(sk)
}

func createKeyBlockProofSk() (root *cell.ProofSkeleton, sk *cell.ProofSkeleton) {
	root = cell.CreateProofSkeleton()
	extraSk := root.ProofRef(3)
	customSk := extraSk.ProofRef(3)
	return root, customSk.ProofRef(3)
}

func createBlockProofSk() *cell.ProofSkeleton {
	sk := cell.CreateProofSkeleton().ProofRef(0)
	return sk
}
