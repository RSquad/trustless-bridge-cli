package txchecker

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/spf13/viper"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type TxCheckerContract struct {
	Addr      *address.Address
	tonClient *tonclient.TonClient
}

type StateInit struct {
	InitData *InitData
	Code     *cell.Cell
}

type InitData struct {
	LiteClientAddr *address.Address
}

func InitDataToCell(initData *InitData) *cell.Cell {
	return cell.BeginCell().
		MustStoreAddr(initData.LiteClientAddr).
		EndCell()
}

func New(
	addr *address.Address,
	tonClient *tonclient.TonClient,
) *TxCheckerContract {
	return &TxCheckerContract{addr, tonClient}
}

func DeployTxChecker(ctx context.Context, tonClient *tonclient.TonClient, initData *InitData) (*address.Address, error) {
	wallet := tonClient.GetWallet()

	msgBody := cell.BeginCell().EndCell()

	codeHex := viper.GetString("tx_checker_code")
	codeBytes, err := hex.DecodeString(codeHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx checker code: %w", err)
	}
	codeCell, err := cell.FromBOC(codeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tx checker code: %w", err)
	}

	addr, _, _, err := wallet.DeployContractWaitTransaction(context.Background(), tlb.MustFromTON("0.1"),
		msgBody,
		codeCell,
		InitDataToCell(initData))

	return addr, err
}
