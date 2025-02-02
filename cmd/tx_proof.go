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

var txProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "Builds transaction proof within a block",
	Long: `This command constructs a proof for a transaction contained within a specified block.
The proof is generated based on the transaction hash and the block file provided.
By default, the proof is output in hexadecimal format.
Usage example: 
    trustless-bridge-cli tx proof -t <transaction_hash> -b <path_to_block.boc>
You can specify the output format using the -f flag, with options 'hex' for hexadecimal or 'bin' for binary.`,
	Run: runTxProof,
}

func init() {
	txCmd.AddCommand(txProofCmd)
	txProofCmd.Flags().BytesHexP("tx-hash", "t", nil, "Transaction hash in hexadecimal format")
	txProofCmd.Flags().StringP("block-boc-path", "b", "", "Path to the BOC file containing the block")
	txProofCmd.MarkFlagRequired("tx-hash")
	txProofCmd.MarkFlagRequired("block-boc-path")
	txProofCmd.Flags().StringP("output-format", "f", "hex", "Output format options: 'bin' for binary, 'hex' for hexadecimal")
}

func runTxProof(cmd *cobra.Command, args []string) {
	txHash, err := cmd.Flags().GetBytesHex("tx-hash")
	if err != nil {
		panic(err)
	}
	blockBocPath, err := cmd.Flags().GetString("block-boc-path")
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

	txProofCell, _, err := txutils.BuildTxProof(blockCell, txHash)
	if err != nil {
		panic(err)
	}

	outputFormattedProof(cmd, txProofCell)
}

func outputFormattedProof(cmd *cobra.Command, proofCell *cell.Cell) {
	outputFormat, err := cmd.Flags().GetString("output-format")
	if err != nil {
		panic(err)
	}

	switch outputFormat {
	case "hex":
		fmt.Printf("%x\n", proofCell.ToBOC())
	case "bin":
		fallthrough
	default:
		os.Stdout.Write(proofCell.ToBOC())
	}
}
