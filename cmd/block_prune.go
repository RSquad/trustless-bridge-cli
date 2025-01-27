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

	"github.com/rsquad/trustless-bridge-cli/internal/blockutils"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

var blockPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune a block to remove unnecessary data",
	Run:   runBlockPrune,
}

func init() {
	blockCmd.AddCommand(blockPruneCmd)
	blockPruneCmd.Flags().StringP("input-file", "i", "", "Input file")
	blockPruneCmd.Flags().BoolP("include-proof-header", "h", false, "Include proof header")
	blockPruneCmd.Flags().StringP("output-format", "f", "bin", "Output format: bin, hex")
	blockPruneCmd.MarkFlagRequired("input-file")
}

func runBlockPrune(cmd *cobra.Command, args []string) {
	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		panic(err)
	}
	outputFormat, err := cmd.Flags().GetString("output-format")
	if err != nil {
		panic(err)
	}
	includeProofHeader, err := cmd.Flags().GetBool("include-proof-header")
	if err != nil {
		panic(err)
	}

	blockBOC, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var result *cell.Cell
	if includeProofHeader {
		result, err = blockutils.BuildBlockProof(blockBOC)
		if err != nil {
			panic(err)
		}
	} else {
		result, err = blockutils.PruneBlock(blockBOC)
		if err != nil {
			panic(err)
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
