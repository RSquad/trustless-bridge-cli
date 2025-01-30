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
	"fmt"
	"os"

	"github.com/rsquad/trustless-bridge-cli/internal/txutils"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

// txProofCmd represents the proofTx command
var txProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "build transaction proof from block",
	Long: `Define transaction hash and path to block BOC file.
	Examples: 
    cli tx proof -t A43DEBA96E0815151645411FEF0FE7E54FF35500B02310A19E7A89AFDFA58194 -b EFFC17EF8FE824E6A039944F3BFC9CEC4A5D9F74D8D93122243EDBD7BF5D4123.boc`,
	Run: runTxProof,
}

func init() {
	txCmd.AddCommand(txProofCmd)
	txProofCmd.Flags().BytesHexP("tx-hash", "t", nil, "tx hash")
	txProofCmd.Flags().StringP("block-boc", "b", "", "path to boc file with block")
	txProofCmd.MarkFlagRequired("tx-hash")
	txProofCmd.MarkFlagRequired("block-boc")
}

func runTxProof(cmd *cobra.Command, args []string) {
	blockBocPath, err := cmd.Flags().GetString("block-boc")
	if err != nil {
		panic(err)
	}
	txHash, err := cmd.Flags().GetBytesHex("tx-hash")
	if err != nil {
		panic(err)
	}

	blockBOC, err := os.ReadFile(blockBocPath)
	if err != nil {
		panic(err)
	}
	blockCell, err := cell.FromBOC(blockBOC)
	if err != nil {
		panic(err)
	}
	txProofCell, err := txutils.BuildTxProof(blockCell, txHash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("tx proof: %s\n\n\n", txProofCell.Dump())
}
