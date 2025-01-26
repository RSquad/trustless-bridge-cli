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

var blockProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "Generate a proof for a block",
	Long: `Generate a proof for a block from the specified input file.
This command processes the block data and outputs the proof in the desired format.
Supported output formats are binary and hexadecimal.`,
	Run: runBlockProof,
}

func init() {
	blockCmd.AddCommand(blockProofCmd)
	blockProofCmd.Flags().StringP("input-file", "i", "", "Input file")
	blockProofCmd.Flags().StringP("output-format", "f", "bin", "Output format: bin, hex")
	blockProofCmd.MarkFlagRequired("input-file")
}

func runBlockProof(cmd *cobra.Command, args []string) {
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

	blockProof, err := blockutils.BuildBlockProof(blockBOC)
	if err != nil {
		panic(err)
	}

	switch outputFormat {
	case "hex":
		fmt.Printf("%x\n", blockProof.ToBOC())

	case "bin":
		fallthrough
	default:
		os.Stdout.Write(blockProof.ToBOC())
	}
}
