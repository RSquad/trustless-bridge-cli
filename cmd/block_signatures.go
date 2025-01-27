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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/adnl"
	"github.com/xssnick/tonutils-go/tl"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

var blockSignaturesCmd = &cobra.Command{
	Use:   "signatures",
	Short: "This command extracts and returns block signatures for a specified block in the masterchain",
	Long: `This command retrieves block proofs from the TON network, locates the validator signatures associated with a given block, and outputs them in one of three formats:
-	json: An array of objects like [{ "pubkey": "...", "signature": "..." }].
-	bin: A BOC containing a TLB-encoded dictionary of type Dict<int256, 512> that maps each validator's public key (int256) to the 512-bit signature.
-	hex: The same BOC as in bin mode, but presented as a hex-encoded string.`,
	Run: runBlockSignatures,
}

func init() {
	blockCmd.AddCommand(blockSignaturesCmd)
	blockSignaturesCmd.Flags().Uint32P("seqno", "s", 0, "Block seqno")
	blockSignaturesCmd.Flags().StringP("output-format", "f", "hex", "Output format: json, bin, hex")
	blockSignaturesCmd.MarkFlagRequired("seqno")
}

func runBlockSignatures(cmd *cobra.Command, args []string) {
	outputFormat, err := cmd.Flags().GetString("output-format")
	if err != nil {
		panic(err)
	}
	seqno, err := cmd.Flags().GetUint32("seqno")
	if err != nil {
		panic(err)
	}
	workchain := int32(-1)

	blockIDExt, err := tonClient.API.LookupBlock(context.Background(), workchain, 0, seqno)
	if err != nil {
		panic(err)
	}
	block, err := tonClient.API.GetBlockData(context.Background(), blockIDExt)
	if err != nil {
		panic(err)
	}

	prevKeyBlockIDExt, err := tonClient.API.LookupBlock(
		context.Background(),
		block.BlockInfo.Shard.WorkchainID,
		int64(block.BlockInfo.Shard.GetShardID()),
		block.BlockInfo.PrevKeyBlockSeqno,
	)
	if err != nil {
		panic(err)
	}
	prevKeyBlock, err := tonClient.API.GetBlockData(context.Background(), prevKeyBlockIDExt)
	if err != nil {
		panic(err)
	}

	validatorsRootCell, err := prevKeyBlock.Extra.Custom.ConfigParams.Config.Params.LoadValueByIntKey(big.NewInt(34))
	if err != nil {
		panic(err)
	}
	var validatorSet tlb.ValidatorSetAny
	if err = tlb.LoadFromCell(&validatorSet, validatorsRootCell.MustLoadRef()); err != nil {
		panic(err)
	}

	validators, err := tonClient.GetMainValidators(validatorSet)
	if err != nil {
		panic(err)
	}

	blockProof, err := tonClient.API.GetBlockProof(
		context.Background(),
		prevKeyBlockIDExt,
		blockIDExt,
	)
	if err != nil {
		panic(err)
	}

	signatures := extractSignatures(blockProof)

	signaturesMap, err := mapValidatorsToSignatures(validators, signatures)
	if err != nil {
		panic(err)
	}

	switch outputFormat {
	case "json":
		stringKeyMap := make(map[string]string)
		for key, value := range signaturesMap {
			stringKeyMap[hex.EncodeToString(key[:])] = hex.EncodeToString(value)
		}

		jsonData, err := json.MarshalIndent(stringKeyMap, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonData))
	case "hex":
		dict := cell.NewDict(256)
		for key, value := range signaturesMap {
			keyCell := cell.BeginCell().MustStoreSlice(key[:], 256).EndCell()
			valueCell := cell.BeginCell().MustStoreSlice(value, 512).EndCell()
			dict.Set(keyCell, valueCell)
		}
		fmt.Printf("%x\n", dict.AsCell().ToBOC())
	case "bin":
		fallthrough
	default:
		dict := cell.NewDict(256)
		for key, value := range signaturesMap {
			keyCell := cell.BeginCell().MustStoreSlice(key[:], 256).EndCell()
			valueCell := cell.BeginCell().MustStoreSlice(value, 512).EndCell()
			dict.Set(keyCell, valueCell)
		}
		os.Stdout.Write(dict.AsCell().ToBOC())
	}
}

func extractSignatures(proof *ton.PartialBlockProof) []ton.Signature {
	var signatures []ton.Signature
	for _, step := range proof.Steps {
		if fwd, ok := step.(ton.BlockLinkForward); ok {
			signatures = append(signatures, fwd.SignatureSet.Signatures...)
		}
	}
	return signatures
}

func mapValidatorsToSignatures(
	validators []*tlb.ValidatorAddr,
	signatures []ton.Signature,
) (map[[32]byte][]byte, error) {
	result := map[[32]byte][]byte{}
	validatorsMap := map[string]*tlb.ValidatorAddr{}
	for _, v := range validators {
		kid, err := tl.Hash(adnl.PublicKeyED25519{Key: v.PublicKey.Key})
		if err != nil {
			return nil, err
		}
		validatorsMap[string(kid)] = v
	}

	for _, signature := range signatures {
		validator := validatorsMap[string(signature.NodeIDShort)]
		var key [32]byte
		copy(key[:], validator.PublicKey.Key)
		result[key] = signature.Signature
	}

	return result, nil
}
