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
)

var blockPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune unnecessary data from a blockchain block",
	Long: `Prune a block to remove unnecessary data, reducing its size.
This command processes the block data from the specified input file and outputs
the pruned block in the desired format. Supported output formats are binary and hexadecimal.`,
	Run: runBlockPrune,
}

func init() {
	blockCmd.AddCommand(blockPruneCmd)
	blockPruneCmd.Flags().StringP("input-file", "i", "", "Input file")
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

	blockBOC, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	prunedBlock, err := blockutils.PruneBlock(blockBOC)
	if err != nil {
		panic(err)
	}

	switch outputFormat {
	case "hex":
		fmt.Printf("%x\n", prunedBlock.ToBOC())

	case "bin":
		fallthrough
	default:
		os.Stdout.Write(prunedBlock.ToBOC())
	}
}
