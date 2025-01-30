package main

import (
	"math/big"
	"testing"

	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type ExampleStruct struct {
	A     uint32           `tlb:"## 24"`
	DictA *cell.Dictionary `tlb:"dict 32"`
}

func prepareExampleCell(keys []int64) (*cell.Cell, *ExampleStruct) {
	data := ExampleStruct{
		A:     0xAABBCC,
		DictA: cell.NewDict(32),
	}

	valueData := cell.BeginCell().
		MustStoreUInt(0x11223311, 32).
		MustStoreRef(cell.BeginCell().
			MustStoreStringSnake("hello tonutils-go").
			MustStoreRef(cell.BeginCell().MustStoreUInt(0xFFFFFFFF, 32).MustStoreUInt(0xAAAAAAAA, 32).EndCell()).
			EndCell()).
		EndCell()

	for _, k := range keys {
		data.DictA.SetIntKey(big.NewInt(k), valueData)
	}

	cl, err := tlb.ToCell(data)
	if err != nil {
		panic(err)
	}

	return cl, &data
}

func TestCreateDict3Proof(t *testing.T) {
	keys := []int64{777, 778, 1}
	exampleCell, data := prepareExampleCell(keys)

	println(exampleCell.Dump())

	sk := cell.CreateProofSkeleton()
	skDictA := sk.ProofRef(0)
	key := cell.BeginCell().MustStoreUInt(778, 32).EndCell()
	_, skKey, err := data.DictA.LoadValueWithProof(key, skDictA)
	if err != nil {
		panic(err)
	}
	skKey.SetRecursive()

	proof, err := exampleCell.CreateProof(sk)
	if err != nil {
		panic(err)
	}

	println("PROOF\n", proof.Dump())
}

func TestCreateDict1Proof(t *testing.T) {
	keys := []int64{778}
	exampleCell, data := prepareExampleCell(keys)

	println("Print Struct\n", exampleCell.Dump())

	sk := cell.CreateProofSkeleton()
	skDictA := sk.ProofRef(0)
	key := cell.BeginCell().MustStoreUInt(778, 32).EndCell()
	_, skKey, err := data.DictA.LoadValueWithProof(key, skDictA)
	if err != nil {
		panic(err)
	}
	skKey.SetRecursive()

	proof, err := exampleCell.CreateProof(sk)
	if err != nil {
		panic(err)
	}

	println("PROOF 1 key\n", proof.Dump())
}
