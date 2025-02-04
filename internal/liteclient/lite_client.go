package liteclient

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/spf13/viper"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	opCodeNewKeyBlock       = 0x11a78ffe
	opCodeNewKeyBlockAnswer = 0xff8ff4e1
	opCodeCheckBlock        = 0x8eaa9d76
	opCodeCheckBlockAnswer  = 0xce02b807
)

type LiteClientContract struct {
	Addr      *address.Address
	tonClient *tonclient.TonClient
}

type StateInit struct {
	InitData *InitData
	Code     *cell.Cell
}

type InitData struct {
	EpochHash             []byte
	ValidatorsTotalWeight uint64
	ValidatorDict         *cell.Dictionary
}

func InitDataToCell(initData *InitData) *cell.Cell {
	return cell.BeginCell().
		MustStoreUInt(initData.ValidatorsTotalWeight, 64).
		MustStoreSlice(initData.EpochHash, 256).
		MustStoreDict(initData.ValidatorDict).
		EndCell()
}

func New(
	addr *address.Address,
	tonClient *tonclient.TonClient,
) *LiteClientContract {
	return &LiteClientContract{addr, tonClient}
}

func (c *LiteClientContract) SendNewKeyBlock(
	ctx context.Context,
	fileHash []byte,
	blockProofCell *cell.Cell,
	signaturesDict *cell.Dictionary,
) (*tlb.Transaction, *ton.BlockIDExt, error) {
	w := c.tonClient.GetWallet()

	payload := cell.BeginCell().
		MustStoreUInt(opCodeNewKeyBlock, 32).
		MustStoreUInt(0, 64).
		MustStoreRef(cell.BeginCell().
			MustStoreSlice(fileHash, 256).
			MustStoreRef(blockProofCell).
			EndCell()).
		MustStoreDict(signaturesDict).
		EndCell()

	message := wallet.SimpleMessage(c.Addr, tlb.MustFromTON("0.2"), payload)

	return w.SendWaitTransaction(ctx, message)
}

func (c *LiteClientContract) SendCheckBlock(
	ctx context.Context,
	fileHash []byte,
	blockProofCell *cell.Cell,
	signaturesDict *cell.Dictionary,
) (*tlb.Transaction, *ton.BlockIDExt, error) {
	w := c.tonClient.GetWallet()

	payload := cell.BeginCell().
		MustStoreUInt(opCodeCheckBlock, 32).
		MustStoreUInt(0, 64).
		MustStoreRef(
			cell.BeginCell().
				MustStoreSlice(fileHash, 256).
				MustStoreRef(blockProofCell).
				EndCell(),
		).
		MustStoreDict(signaturesDict).
		EndCell()

	message := wallet.SimpleMessage(c.Addr, tlb.MustFromTON("0.2"), payload)

	return w.SendWaitTransaction(ctx, message)
}

func DeployLiteClient(ctx context.Context, tonClient *tonclient.TonClient, initData *InitData) (*address.Address, error) {
	wallet := tonClient.GetWallet()

	msgBody := cell.BeginCell().EndCell()

	codeHex := viper.GetString("lite_client_code")
	codeBytes, err := hex.DecodeString(codeHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode lite client code: %w", err)
	}
	codeCell, err := cell.FromBOC(codeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lite client code: %w", err)
	}

	addr, _, _, err := tonclient.DeployContractWaitTransaction(
		context.Background(),
		wallet,
		255, // workchain
		tlb.MustFromTON("0.2"),
		msgBody,
		codeCell,
		InitDataToCell(initData),
	)

	return addr, err
}
