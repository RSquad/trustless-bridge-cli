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

	"github.com/spf13/cobra"
)

var blockProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "Generate a proof for a block",
	Run:   runBlockProof,
}

func init() {
	blockCmd.AddCommand(blockProofCmd)
}

func runBlockProof(cmd *cobra.Command, args []string) {
	fmt.Println("not implemented")
}
