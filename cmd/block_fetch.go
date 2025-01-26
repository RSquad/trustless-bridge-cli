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
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var blockFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a block from the blockchain using its seqno",
	Long: `Fetch a block from the blockchain by specifying the seqno and workchain.
You can choose the output format as JSON, binary, or hexadecimal.`,
	Run: run,
}

func init() {
	blockCmd.AddCommand(blockFetchCmd)
	blockFetchCmd.Flags().Uint32P("seqno", "s", 0, "Block seqno")
	blockFetchCmd.Flags().Int32P("workchain", "w", -1, "Workchain")
	blockFetchCmd.Flags().StringP("output-format", "f", "hex", "Output format: json, bin, hex")
	blockFetchCmd.MarkFlagRequired("seqno")
}

func run(cmd *cobra.Command, args []string) {
	outputFormat, err := cmd.Flags().GetString("output-format")
	if err != nil {
		panic(err)
	}
	workchain, err := cmd.Flags().GetInt32("workchain")
	if err != nil {
		panic(err)
	}
	seqno, err := cmd.Flags().GetUint32("seqno")
	if err != nil {
		panic(err)
	}

	blockIDExt, err := tonClient.API.LookupBlock(context.Background(), workchain, 0, seqno)
	if err != nil {
		panic(err)
	}

	switch outputFormat {
	case "json":
		block, err := tonClient.API.GetBlockData(context.Background(), blockIDExt)
		if err != nil {
			panic(err)
		}
		blockJSON, err := json.MarshalIndent(block, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", blockJSON)

	case "hex":
		blockBOC, err := tonClient.GetBlockBOC(context.Background(), blockIDExt)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%x\n", blockBOC)

	case "bin":
		fallthrough
	default:
		blockBOC, err := tonClient.GetBlockBOC(context.Background(), blockIDExt)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(blockBOC)
	}
}
