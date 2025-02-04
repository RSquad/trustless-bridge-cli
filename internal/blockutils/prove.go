package blockutils

import (
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func BuildBlockProof(blockBOC []byte) (*cell.Cell, error) {
	blockCell, err := cell.FromBOC(blockBOC)
	if err != nil {
		return nil, err
	}

	var block tlb.Block
	err = tlb.LoadFromCell(&block, blockCell.BeginParse())
	if err != nil {
		return nil, err
	}

	if block.Extra == nil || block.Extra.Custom == nil || block.Extra.Custom.ConfigParams == nil {
		rootSk := createBlockProofSk()
		return blockCell.CreateProof(rootSk)
	}

	configRefIndex := 3
	if block.Extra.Custom.ShardHashes.IsEmpty() {
		configRefIndex -= 1
	}
	if block.Extra.Custom.ShardFees.IsEmpty() {
		configRefIndex -= 1
	}
	rootSk, configSk := createKeyBlockProofSk(configRefIndex)
	_, config34Sk, err := block.Extra.Custom.ConfigParams.Config.Params.LoadValueWithProof(
		cell.BeginCell().MustStoreUInt(34, 32).EndCell(),
		configSk,
	)
	if err != nil {
		return nil, err
	}
	config34Sk.SetRecursive()

	return blockCell.CreateProof(rootSk)
}

func createKeyBlockProofSk(configIdx int) (rootSk *cell.ProofSkeleton, configSk *cell.ProofSkeleton) {
	rootSk = cell.CreateProofSkeleton()
	extraSk := rootSk.ProofRef(3)
	customSk := extraSk.ProofRef(3)
	configSk = customSk.ProofRef(configIdx)
	return rootSk, configSk
}

func createBlockProofSk() (rootSk *cell.ProofSkeleton) {
	rootSk = cell.CreateProofSkeleton()
	rootSk.ProofRef(0)
	return rootSk
}
