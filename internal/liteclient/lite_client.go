package liteclient

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

type LiteClientContract struct {
	Addr      *address.Address
	tonClient *tonclient.TonClient
	ctx       context.Context
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
	ctx context.Context,
) *LiteClientContract {
	return &LiteClientContract{addr, tonClient, ctx}
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

	addr, _, _, err := wallet.DeployContractWaitTransaction(context.Background(), tlb.MustFromTON("0.1"),
		msgBody,
		codeCell,
		InitDataToCell(initData))

	return addr, err
}
