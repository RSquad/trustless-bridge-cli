package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/rsquad/trustless-bridge-cli/internal/blockutils"
	"github.com/rsquad/trustless-bridge-cli/internal/liteclient"
	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/address"
)

var sendNewKeyBlockCmd = &cobra.Command{
	Use:   "new-key-block",
	Short: "Send a new_key_block message to a LiteClient",
	Long: `This command sends a new_key_block message to a LiteClient.
If the network is specified as testnet, the system will fetch a block from fastnet
and send a new_key_block message to LiteClient in testnet.`,
	RunE: runSendNewKeyBlock,
}

func init() {
	sendCmd.AddCommand(sendNewKeyBlockCmd)
	sendNewKeyBlockCmd.Flags().Uint32P("seqno", "s", 0, "Block seqno")
	sendNewKeyBlockCmd.MarkFlagRequired("seqno")
}

func runSendNewKeyBlock(cmd *cobra.Command, args []string) error {
	network, err := cmd.Flags().GetString("network")
	if err != nil {
		network = "testnet"
	}
	seqno, err := cmd.Flags().GetUint32("seqno")
	if err != nil {
		return fmt.Errorf("failed to get seqno: %w", err)
	}
	addrStr, err := cmd.Flags().GetString("address")
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}
	addr, err := address.ParseAddr(addrStr)
	if err != nil {
		return fmt.Errorf("failed to parse address: %w", err)
	}

	oppositeNetwork := "testnet"
	if network == "testnet" {
		oppositeNetwork = "fastnet"
	}
	oppositeTonClient, err := tonclient.NewTonClientNetwork(oppositeNetwork)
	if err != nil {
		return fmt.Errorf("failed to create TonClient: %w", err)
	}

	log.Printf("Attention: You are sending a message to the %s network with block %d from %s network", network, seqno, oppositeNetwork)

	blockIDExt, blockBOC, err := blockutils.FetchMasterchainBlockBOC(context.Background(), oppositeTonClient, seqno)
	if err != nil {
		return fmt.Errorf("failed to fetch masterchain block: %w", err)
	}

	liteClient := liteclient.New(addr, tonClient)

	signaturesMap, err := GetBlockSignatures(seqno, oppositeTonClient)
	if err != nil {
		return fmt.Errorf("failed to get block signatures: %w", err)
	}
	signaturesDict := SignaturesMapToDict(signaturesMap)

	blockProof, err := blockutils.BuildBlockProof(blockBOC)
	if err != nil {
		return fmt.Errorf("failed to build block proof: %w", err)
	}

	sendTx, blockIDExt, err := liteClient.SendNewKeyBlock(
		context.Background(),
		blockIDExt.FileHash,
		blockProof,
		signaturesDict,
	)
	if err != nil {
		return fmt.Errorf("failed to send new key block: %w", err)
	}

	fmt.Printf("NewKeyBlock successfully sent\n")
	fmt.Printf("With transaction lt: %v, hash: %x\n", sendTx.LT, sendTx.Hash)
	fmt.Printf("In block: %v\n", blockIDExt.SeqNo)

	return nil
}
