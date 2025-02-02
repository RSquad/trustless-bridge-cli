package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/rsquad/trustless-bridge-cli/internal/blockutils"
	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/rsquad/trustless-bridge-cli/internal/txchecker"
	"github.com/rsquad/trustless-bridge-cli/internal/txutils"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/address"
)

var sendCheckTxCmd = &cobra.Command{
	Use:   "check-tx",
	Short: "Send a check_transaction message to a TxChecker",
	Long: `This command sends a check_transaction message to a TxChecker.
If the network is specified as testnet, the system will fetch a block and transaction from fastnet
and send a check_transaction message to TxChecker in testnet.`,
	RunE: runSendCheckTx,
}

func init() {
	sendCmd.AddCommand(sendCheckTxCmd)
	sendCheckTxCmd.Flags().Uint32P("seqno", "s", 0, "Block seqno")
	sendCheckTxCmd.Flags().BytesHexP("tx-hash", "t", nil, "Transaction hash in hexadecimal format")
	sendCheckTxCmd.MarkFlagRequired("seqno")
	sendCheckTxCmd.MarkFlagRequired("tx-hash")
}

func runSendCheckTx(cmd *cobra.Command, args []string) error {
	network, err := cmd.Flags().GetString("network")
	if err != nil {
		network = "testnet"
	}
	seqno, err := cmd.Flags().GetUint32("seqno")
	if err != nil {
		return fmt.Errorf("failed to get seqno: %w", err)
	}
	txHash, err := cmd.Flags().GetBytesHex("tx-hash")
	if err != nil {
		return fmt.Errorf("failed to get tx hash: %w", err)
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

	log.Printf("Attention: You are sending a message to the %s network with transaction %x and block %d from %s network", network, txHash, seqno, oppositeNetwork)

	blockCell, err := blockutils.FetchMasterchainBlockCell(context.Background(), oppositeTonClient, seqno)
	if err != nil {
		return fmt.Errorf("failed to fetch masterchain block: %w", err)
	}

	txProofCell, tx, err := txutils.BuildTxProof(blockCell, txHash)
	if err != nil {
		return fmt.Errorf("failed to build tx proof: %w", err)
	}

	txChecker := txchecker.New(addr, tonClient)

	tx.Hash = txHash

	sendTx, blockIDExt, err := txChecker.SendCheckTx(
		context.Background(),
		txchecker.TxToCell(tx),
		txProofCell,
		blockCell,
	)
	if err != nil {
		return fmt.Errorf("failed to send check tx: %w", err)
	}

	fmt.Printf("CheckTx for tx %x successfully sent\n", txHash)
	fmt.Printf("With transaction lt: %v, hash: %x\n", sendTx.LT, sendTx.Hash)
	fmt.Printf("In block: %v\n", blockIDExt.SeqNo)

	return nil
}
