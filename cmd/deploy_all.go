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
	"fmt"

	"github.com/rsquad/trustless-bridge-cli/internal/blockutils"
	"github.com/rsquad/trustless-bridge-cli/internal/liteclient"
	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/rsquad/trustless-bridge-cli/internal/txchecker"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

var deployAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Deploy system contracts",
	Long: `This command deploys system contracts to the opposite network.
If the network is specified as testnet, the system will fetch a block from fastnet
and deploy the system to testnet using the block from fastnet, and vice versa.`,
	RunE: runDeployAll,
}

func init() {
	deployCmd.AddCommand(deployAllCmd)
	deployAllCmd.Flags().Uint32P("trusted-block-seqno", "s", 0, "Trusted block seqno")
	deployAllCmd.Flags().Int8P("workchain", "w", 0, "Workchain")
	deployAllCmd.MarkFlagRequired("trusted-block-seqno")
	deployAllCmd.MarkFlagRequired("workchain")
}

func runDeployAll(cmd *cobra.Command, args []string) error {
	trustedBlockSeqno, err := cmd.Flags().GetUint32("trusted-block-seqno")
	if err != nil {
		return fmt.Errorf("failed to get trusted block seqno: %w", err)
	}
	network, err := cmd.Flags().GetString("network")
	if err != nil {
		network = "testnet"
	}
	oppositeNetwork := "testnet"
	if network == "testnet" {
		oppositeNetwork = "fastnet"
	}
	oppositeTonClient, err := tonclient.NewTonClientNetwork(oppositeNetwork)
	if err != nil {
		return fmt.Errorf("failed to create TonClient: %w", err)
	}
	wc, err := cmd.Flags().GetInt8("workchain")
	if err != nil {
		return fmt.Errorf("failed to get workchain: %w", err)
	}
	wcb := byte(wc)
	if wc != 0 {
		wcb = 255
	}

	fmt.Printf("Attention: You are deploying contracts to the %s network with block %d from %s network\n", network, trustedBlockSeqno, oppositeNetwork)

	trustedBlock, err := blockutils.FetchMasterchainBlock(context.Background(), oppositeTonClient, trustedBlockSeqno)
	if err != nil {
		return fmt.Errorf("failed to fetch masterchain block: %w", err)
	}

	if !trustedBlock.BlockInfo.KeyBlock {
		fmt.Printf("given trusted block is not a key block: %v\n", trustedBlock.BlockInfo.SeqNo)
		fmt.Printf("switch to last key block with seqno: %v\n", trustedBlock.BlockInfo.PrevKeyBlockSeqno)
		trustedBlockSeqno = trustedBlock.BlockInfo.PrevKeyBlockSeqno

		trustedBlock, err = blockutils.FetchMasterchainBlock(context.Background(), oppositeTonClient, trustedBlockSeqno)
		if err != nil {
			return fmt.Errorf("failed to fetch masterchain block: %w", err)
		}
	}

	validators, validatorsTotalWeight, epochHash, err := blockutils.ExtractMainValidators(trustedBlock, oppositeTonClient)
	if err != nil {
		return fmt.Errorf("failed to extract main validators: %w", err)
	}

	validatorDict := cell.NewDict(256)

	for _, validator := range validators {
		validatorDict.Set(
			cell.BeginCell().MustStoreSlice(validator.PublicKey.Key, 256).EndCell(),
			cell.BeginCell().MustStoreUInt(validator.Weight, 64).EndCell(),
		)
	}

	liteClientAddr, err := liteclient.DeployLiteClient(
		context.Background(),
		tonClient,
		wcb,
		&liteclient.InitData{
			EpochHash:             epochHash,
			ValidatorsTotalWeight: validatorsTotalWeight,
			ValidatorDict:         validatorDict,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to deploy lite client: %w", err)
	}

	txCheckerAddr, err := txchecker.DeployTxChecker(
		context.Background(),
		tonClient,
		wcb,
		&txchecker.InitData{LiteClientAddr: liteClientAddr},
	)
	if err != nil {
		return fmt.Errorf("failed to deploy tx checker: %w", err)
	}

	fmt.Printf("LiteClient successfully deployed: %v\n", liteClientAddr)
	fmt.Printf("TxChecker successfully deployed: %v\n", txCheckerAddr)

	return nil
}
