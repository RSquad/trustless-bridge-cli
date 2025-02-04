package txchecker

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/rsquad/trustless-bridge-cli/internal/wallet"
	"github.com/spf13/viper"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	opCodeCheckTx = 0x91d555f7
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

func (c *TxCheckerContract) SendCheckTx(
	ctx context.Context,
	txCell *cell.Cell,
	proofCell *cell.Cell,
	blockCell *cell.Cell,
) (*tlb.Transaction, *ton.BlockIDExt, error) {
	w := c.tonClient.GetWallet()

	payload := cell.BeginCell().
		MustStoreUInt(opCodeCheckTx, 32).
		MustStoreRef(txCell).
		MustStoreRef(proofCell).
		MustStoreRef(blockCell).
		EndCell()

	message := wallet.SimpleMessage(c.Addr, tlb.MustFromTON("1"), payload)

	return w.SendWaitTransaction(ctx, message)
}

func DeployTxChecker(ctx context.Context, tonClient *tonclient.TonClient, wc byte, initData *InitData) (*address.Address, error) {
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

	addr, _, _, err := tonclient.DeployContractWaitTransaction(
		context.Background(),
		wallet,
		wc,
		tlb.MustFromTON("0.2"),
		msgBody,
		codeCell,
		InitDataToCell(initData),
	)

	return addr, err
}

func TxToCell(tx *tlb.Transaction) *cell.Cell {
	return cell.BeginCell().
		MustStoreSlice(tx.Hash, 256).
		MustStoreSlice(tx.AccountAddr, 256).
		MustStoreUInt(tx.LT, 64).
		EndCell()
}
