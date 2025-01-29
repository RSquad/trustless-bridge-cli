/*
Copyright © 2025 RSquad <hello@rsquad.io>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bytes"
	"fmt"
	"math/big"
	"os"

	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

// txProofCmd represents the proofTx command
var txProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runTxProof,
}

func init() {
	txCmd.AddCommand(txProofCmd)
	txProofCmd.Flags().BytesHexP("account-id", "a", nil, "account addr (uint256)")
	txProofCmd.Flags().BytesHexP("tx-hash", "t", nil, "tx hash")
	txProofCmd.Flags().StringP("block-boc", "b", "", "path to boc file with block")
	txProofCmd.MarkFlagRequired("account-id")
	txProofCmd.MarkFlagRequired("tx-hash")
	txProofCmd.MarkFlagRequired("block-boc")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txProofCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txProofCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runTxProof(cmd *cobra.Command, args []string) {
	blockBocPath, err := cmd.Flags().GetString("block-boc")
	if err != nil {
		panic(err)
	}
	accountId, err := cmd.Flags().GetBytesHex("account-id")
	if err != nil {
		panic(err)
	}
	txHash, err := cmd.Flags().GetBytesHex("tx-hash")
	if err != nil {
		panic(err)
	}
	fmt.Printf("tx hash   : %x\n", txHash)

	blockBOC, err := os.ReadFile(blockBocPath)
	if err != nil {
		panic(err)
	}
	blockCell, err := cell.FromBOC(blockBOC)
	if err != nil {
		panic(err)
	}
	sk := cell.CreateProofSkeleton()
	accountBlocksSk := sk.ProofRef(3).ProofRef(2).ProofRef(0)

	extra := blockCell.MustPeekRef(3)
	accountBlocksCell := extra.MustPeekRef(2)
	accountBlocksDict := accountBlocksCell.BeginParse().MustLoadDict(256)
	accountBlock, accountBlockSk, err := accountBlocksDict.LoadValueWithProof(
		cell.BeginCell().MustStoreBigUInt(new(big.Int).SetBytes(accountId), 256).EndCell(), accountBlocksSk)
	if err != nil {
		panic(err)
	}
	// skip extra value: CurrencyCollection
	accountBlock.MustLoadCoins()
	accountBlock.MustLoadDict(32)
	accountTransTag := accountBlock.MustLoadUInt(4)
	fmt.Println("acc_trans tag", accountTransTag)
	accountBlock.MustLoadBigUInt(256)
	transDict, err := accountBlock.ToDict(64)
	if err != nil {
		panic(err)
	}
	txLT := uint64(0)
	dictItems, err := transDict.LoadAll()
	if err != nil {
		panic(err)
	}
	for _, kv := range dictItems {
		kv.Value.MustLoadCoins()
		kv.Value.MustLoadDict(32)
		hash := kv.Value.MustLoadRef().MustToCell().Hash()
		if bytes.Equal(txHash, hash) {
			txLT = kv.Key.MustLoadUInt(64)
			break
		}
	}
	fmt.Println("txLT", txLT)
	value, _, err := transDict.LoadValueWithProof(cell.BeginCell().MustStoreUInt(txLT, 64).EndCell(), accountBlockSk)
	if err != nil {
		panic(err)
	}

	// skip extra: CurrencyCollection
	value.MustLoadCoins()
	value.MustLoadDict(32)
	txSlice := value.MustLoadRef()
	foundTxHash := txSlice.MustToCell().Hash()
	fmt.Printf("found tx hash: %x\n", foundTxHash)

	txProof, err := blockCell.CreateProof(sk)
	if err != nil {
		panic(err)
	}
	fmt.Printf("tx proof: %s\n\n\n", txProof.Dump())
	//fmt.Printf("block: %s\n\n", blockCell.Dump())
}
