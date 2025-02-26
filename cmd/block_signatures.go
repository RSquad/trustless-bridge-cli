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
	"sort"

	"github.com/rsquad/trustless-bridge-cli/internal/blockutils"
	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/adnl"
	"github.com/xssnick/tonutils-go/tl"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

var blockSignaturesCmd = &cobra.Command{
	Use:   "signatures",
	Short: "This command extracts and returns necessary block signatures for verification in the masterchain",
	Long: `This command retrieves block proofs from the TON network, locates the validator signatures associated with a given block, and outputs only the necessary signatures for verification in one of three formats:
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

	signaturesMap, err := GetBlockSignatures(seqno, tonClient)
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
		fmt.Printf("%x\n", SignaturesMapToDict(signaturesMap).AsCell().ToBOC())
	case "bin":
		fallthrough
	default:
		os.Stdout.Write(SignaturesMapToDict(signaturesMap).AsCell().ToBOC())
	}
}

func GetBlockSignatures(seqno uint32, tonClient *tonclient.TonClient) (map[[32]byte][]byte, error) {
	workchain := int32(-1)

	blockIDExt, err := tonClient.API.LookupBlock(context.Background(), workchain, 0, seqno)
	if err != nil {
		return nil, err
	}
	block, err := tonClient.API.GetBlockData(context.Background(), blockIDExt)
	if err != nil {
		return nil, err
	}

	prevBlockIDExt, err := tonClient.API.LookupBlock(
		context.Background(),
		block.BlockInfo.Shard.WorkchainID,
		int64(block.BlockInfo.Shard.GetShardID()),
		seqno-1,
	)
	if err != nil {
		return nil, err
	}

	prevKeyBlockIDExt, err := tonClient.API.LookupBlock(
		context.Background(),
		block.BlockInfo.Shard.WorkchainID,
		int64(block.BlockInfo.Shard.GetShardID()),
		block.BlockInfo.PrevKeyBlockSeqno,
	)
	if err != nil {
		return nil, err
	}
	prevKeyBlock, err := tonClient.API.GetBlockData(context.Background(), prevKeyBlockIDExt)
	if err != nil {
		return nil, err
	}

	validatorsRootCell, err := prevKeyBlock.Extra.Custom.ConfigParams.Config.Params.LoadValueByIntKey(big.NewInt(34))
	if err != nil {
		return nil, err
	}
	var validatorSet tlb.ValidatorSetAny
	if err = tlb.LoadFromCell(&validatorSet, validatorsRootCell.MustLoadRef()); err != nil {
		return nil, err
	}

	validators, _, _, err := blockutils.ExtractMainValidators(prevKeyBlock, tonClient)
	if err != nil {
		return nil, err
	}

	blockProof, err := tonClient.GetBlockProofExt(
		context.Background(),
		prevBlockIDExt,
		blockIDExt,
	)
	if err != nil {
		return nil, err
	}

	signatures := extractSignatures(blockProof)

	signaturesMap, err := mapValidatorsToSignatures(validators, signatures)
	if err != nil {
		return nil, err
	}

	return signaturesMap, nil
}

func extractSignatures(proof *ton.PartialBlockProof) []ton.Signature {
	for i := len(proof.Steps) - 1; i >= 0; i-- {
		if fwd, ok := proof.Steps[i].(ton.BlockLinkForward); ok {
			return fwd.SignatureSet.Signatures
		}
	}
	return nil
}

func mapValidatorsToSignatures(
	validators []*tlb.ValidatorAddr,
	signatures []ton.Signature,
) (map[[32]byte][]byte, error) {
	var totalWeight uint64
	validatorsMap := make(map[string]*tlb.ValidatorAddr)

	for _, validator := range validators {
		kid, err := tl.Hash(adnl.PublicKeyED25519{Key: validator.PublicKey.Key})
		if err != nil {
			return nil, err
		}
		validatorsMap[string(kid)] = validator
		totalWeight += validator.Weight
	}

	type signerInfo struct {
		key       [32]byte
		weight    uint64
		signature []byte
	}
	var signers []signerInfo

	for _, s := range signatures {
		if v, ok := validatorsMap[string(s.NodeIDShort)]; ok {
			var key [32]byte
			copy(key[:], v.PublicKey.Key)
			signers = append(signers, signerInfo{
				key:       key,
				weight:    v.Weight,
				signature: s.Signature,
			})
		}
	}

	sort.Slice(signers, func(i, j int) bool {
		return signers[i].weight > signers[j].weight
	})

	var signedWeight uint64
	minSignatures := make(map[[32]byte][]byte)

	for _, signer := range signers {
		minSignatures[signer.key] = signer.signature
		signedWeight += signer.weight
		if 3*signedWeight > 2*totalWeight {
			break
		}
	}

	if 3*signedWeight <= 2*totalWeight {
		return nil, fmt.Errorf("insufficient signed weight (%d/%d)", 3*signedWeight, 2*totalWeight)
	}

	return minSignatures, nil
}

func SignaturesMapToDict(signaturesMap map[[32]byte][]byte) *cell.Dictionary {
	dict := cell.NewDict(256)
	for key, value := range signaturesMap {
		keyCell := cell.BeginCell().MustStoreSlice(key[:], 256).EndCell()
		valueCell := cell.BeginCell().MustStoreSlice(value, 512).EndCell()
		dict.Set(keyCell, valueCell)
	}
	return dict
}
