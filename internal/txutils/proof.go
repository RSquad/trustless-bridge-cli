package txutils

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/xssnick/tonutils-go/tvm/cell"
)

type AccountBlockTx struct {
	accountId         []byte
	txHash            []byte
	txLogicalTime     uint64
	totalAccounts     int
	totalTransactions int
}

func skipCurrencyCollection(s *cell.Slice) {
	s.MustLoadCoins()
	s.MustLoadDict(32)
}

func findAccountBlockTx(accountBlocksDict *cell.Dictionary, txHash []byte) (*AccountBlockTx, error) {
	accounts, err := accountBlocksDict.LoadAll()
	if err != nil {
		return nil, err
	}
	totalAccounts := len(accounts)
	for _, kv := range accounts {
		accountBlock := kv.Value
		accountId := kv.Key.MustLoadBigUInt(256).Bytes()
		skipCurrencyCollection(accountBlock)
		accountTransTag := accountBlock.MustLoadUInt(4)
		if accountTransTag != 5 {
			return nil, fmt.Errorf("ShardAccountBlock has invalid accountBlock (invalid tag)")
		}
		// skip account_addr
		accountBlock.MustLoadBigUInt(256)
		transDict, err := accountBlock.ToDict(64)
		if err != nil {
			return nil, err
		}
		transactions, err := transDict.LoadAll()
		if err != nil {
			return nil, err
		}
		totalTransactions := len(transactions)
		for _, tran := range transactions {
			skipCurrencyCollection(tran.Value)
			hash := tran.Value.MustLoadRef().MustToCell().Hash()
			if bytes.Equal(txHash, hash) {
				txLogicalTime := tran.Key.MustLoadUInt(64)
				return &AccountBlockTx{
					accountId, txHash, txLogicalTime, totalAccounts, totalTransactions,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("tx not found in transactions HashmapAug")
}

func BuildTxProof(blockCell *cell.Cell, txHash []byte) (*cell.Cell, error) {
	rootProofPath := cell.CreateProofSkeleton()
	sk := rootProofPath.ProofRef(3).ProofRef(2).ProofRef(0)

	extra := blockCell.MustPeekRef(3)
	accountBlocksCell := extra.MustPeekRef(2)
	accountBlocksDict := accountBlocksCell.BeginParse().MustLoadDict(256)
	accountBlockTx, err := findAccountBlockTx(accountBlocksDict, txHash)
	if err != nil {
		return nil, err
	}
	accountBlock, accBlockSk, err := accountBlocksDict.LoadValueWithProof(
		cell.BeginCell().MustStoreBigUInt(new(big.Int).SetBytes(accountBlockTx.accountId), 256).EndCell(),
		sk)
	if err != nil {
		return nil, err
	}
	if accountBlockTx.totalAccounts == 1 {
		accBlockSk = sk
	}

	skipCurrencyCollection(accountBlock)
	// skip tag
	accountBlock.MustLoadUInt(4)
	// skip account_addr
	accountBlock.MustLoadBigUInt(256)
	transDict, _ := accountBlock.ToDict(64)

	_, transactionSk, err := transDict.LoadValueWithProof(cell.BeginCell().MustStoreUInt(accountBlockTx.txLogicalTime, 64).EndCell(), accBlockSk)
	if err != nil {
		return nil, err
	}
	if accountBlockTx.totalTransactions == 1 {
		transactionSk = accBlockSk
	}
	transactionSk.SetRecursive()

	txProof, err := blockCell.CreateProof(rootProofPath)
	if err != nil {
		return nil, err
	}
	return txProof, nil
}
