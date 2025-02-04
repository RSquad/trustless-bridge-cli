/*
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
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

var blockProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "Generate a proof for a block (currently works only with blocks from the masterchain)",
	Run:   runBlockProof,
}

func init() {
	blockCmd.AddCommand(blockProofCmd)
	blockProofCmd.Flags().Uint32("from-seqno", 0, "From block seqno")
	blockProofCmd.Flags().Uint32("to-seqno", 0, "To block seqno")
	blockProofCmd.Flags().StringP("output-format", "f", "bin", "Output format: bin, hex")
	blockProofCmd.MarkFlagRequired("from-seqno")
	blockProofCmd.MarkFlagRequired("to-seqno")
}

func runBlockProof(cmd *cobra.Command, args []string) {
	outputFormat, err := cmd.Flags().GetString("output-format")
	if err != nil {
		panic(err)
	}
	fromSeqno, err := cmd.Flags().GetUint32("from-seqno")
	if err != nil {
		panic(err)
	}
	toSeqno, err := cmd.Flags().GetUint32("to-seqno")
	if err != nil {
		panic(err)
	}
	toWorkchain := int32(-1)
	fromWorkchain := int32(-1)

	fromBlockIDExt, err := tonClient.API.LookupBlock(
		context.Background(),
		fromWorkchain,
		0,
		fromSeqno,
	)
	if err != nil {
		panic(err)
	}
	toBlockIDExt, err := tonClient.API.LookupBlock(
		context.Background(),
		toWorkchain,
		0,
		toSeqno,
	)

	if err != nil {
		panic(err)
	}

	blockProof, err := tonClient.API.GetBlockProof(
		context.Background(),
		fromBlockIDExt,
		toBlockIDExt,
	)
	if err != nil {
		panic(err)
	}

	var result *cell.Cell

	for _, step := range blockProof.Steps {
		if back, ok := step.(ton.BlockLinkBackward); ok {
			boc, err := cell.FromBOC(back.Proof)
			if err != nil {
				panic(err)
			}
			result = boc
		}
	}

	switch outputFormat {
	case "hex":
		fmt.Printf("%x\n", result.ToBOC())

	case "bin":
		fallthrough
	default:
		os.Stdout.Write(result.ToBOC())
	}
}
