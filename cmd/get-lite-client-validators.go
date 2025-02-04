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

	"github.com/rsquad/trustless-bridge-cli/internal/liteclient"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/address"
)

var getLiteClientValidatorsCmd = &cobra.Command{
	Use:   "lite-client-validators",
	Short: "Get the validators of the lite client contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		addrStr, err := cmd.Flags().GetString("address")
		if err != nil {
			return fmt.Errorf("failed to get address: %w", err)
		}
		addr, err := address.ParseAddr(addrStr)
		if err != nil {
			return fmt.Errorf("failed to parse address: %w", err)
		}

		liteClient := liteclient.New(addr, tonClient)
		validatorDict, err := liteClient.GetValidators(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get validators: %w", err)
		}
		json, err := liteclient.ValidatorDictToJSON(validatorDict)
		if err != nil {
			return fmt.Errorf("failed to marshal validator dict: %w", err)
		}
		fmt.Printf("Validator dict: %s\n", json)
		return nil
	},
}

func init() {
	getCmd.AddCommand(getLiteClientValidatorsCmd)
}
