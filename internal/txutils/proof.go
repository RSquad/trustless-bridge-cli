package txutils

import (
	"bytes"
	"fmt"

	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func skipCC(s *cell.Slice) {
	var cc tlb.CurrencyCollection
	tlb.LoadFromCellAsProof(&cc, s)
}

func findTxInAccountBlocks(
	accBlocks tlb.ShardAccountBlocks,
	givenTxHash []byte,
) (*tlb.Transaction, error) {
	accBlocksKV, err := accBlocks.Accounts.LoadAll()
	if err != nil {
		return nil, err
	}

	var tx tlb.Transaction

	for _, accKV := range accBlocksKV {
		accCell := accKV.Value
		skipCC(accCell)

		var accBlock tlb.AccountBlock
		err := tlb.LoadFromCell(&accBlock, accCell)
		if err != nil {
			continue
		}

		txs, err := accBlock.Transactions.LoadAll()
		if err != nil {
			continue
		}

		for _, txKV := range txs {
			txV := txKV.Value
			skipCC(txV)

			txCell := txV.MustLoadRef().MustToCell()
			txHash := txCell.Hash()
			err = tlb.LoadFromCell(&tx, txCell.BeginParse())
			if err != nil {
				continue
			}

			if bytes.Equal(givenTxHash, txHash) {
				return &tx, nil
			}
		}
	}
	return nil, fmt.Errorf("tx not found")
}

func BuildTxProof(blockCell *cell.Cell, txHash []byte) (*cell.Cell, error) {
	rootSk := cell.CreateProofSkeleton()
	sk := rootSk.ProofRef(3).ProofRef(2).ProofRef(0)

	var accBlocks tlb.ShardAccountBlocks
	tlb.LoadFromCell(&accBlocks, blockCell.MustPeekRef(3).MustPeekRef(2).BeginParse())
	tx, err := findTxInAccountBlocks(accBlocks, txHash)
	if err != nil {
		return nil, err
	}

	accCell, accBlockSk, err := accBlocks.Accounts.LoadValueWithProof(
		cell.BeginCell().MustStoreSlice(tx.AccountAddr, 256).EndCell(),
		sk)
	if err != nil {
		return nil, err
	}

	skipCC(accCell)

	var accBlock tlb.AccountBlock
	tlb.LoadFromCell(&accBlock, accCell)

	_, txSk, err := accBlock.Transactions.LoadValueWithProof(
		cell.BeginCell().MustStoreUInt(tx.LT, 64).EndCell(),
		accBlockSk,
	)
	if err != nil {
		return nil, err
	}

	txSk.SetRecursive()

	txProof, err := blockCell.CreateProof(rootSk)
	if err != nil {
		return nil, err
	}
	return txProof, nil
}
