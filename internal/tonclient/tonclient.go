package tonclient

import (
	"context"
	"fmt"
	"sort"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tl"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type TonClient struct {
	connPool *liteclient.ConnectionPool
	API      *ton.APIClient
}

func NewTonClient(configUrl string) (*TonClient, error) {
	connPool := liteclient.NewConnectionPool()

	err := connPool.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return nil, err
	}
	apiWrapped := ton.NewAPIClient(connPool).WithRetry(3)
	api, _ := apiWrapped.(*ton.APIClient)

	return &TonClient{connPool: connPool, API: api}, nil
}

func (tc *TonClient) GetBlockBOC(ctx context.Context, block *ton.BlockIDExt) ([]byte, error) {
	var resp tl.Serializable
	err := tc.API.Client().QueryLiteserver(ctx, ton.GetBlockData{ID: block}, &resp)
	if err != nil {
		return nil, err
	}

	switch t := resp.(type) {
	case ton.BlockData:
		return t.Payload, nil
	case ton.LSError:
		return nil, t
	}
	panic("should not happen")
}

func (tc *TonClient) GetMainValidators(validatorSet tlb.ValidatorSetAny) (
	[]*tlb.ValidatorAddr, error,
) {
	var validatorsNum int
	var validatorsListDict *cell.Dictionary

	var definedWeight *uint64
	switch t := validatorSet.Validators.(type) {
	case tlb.ValidatorSet:
		validatorsNum = int(t.Main)
		validatorsListDict = t.List
	case tlb.ValidatorSetExt:
		definedWeight = &t.TotalWeight
		validatorsNum = int(t.Main)
		validatorsListDict = t.List
	default:
		return nil, fmt.Errorf("unknown validator set type")
	}

	type validatorWithKey struct {
		addr *tlb.ValidatorAddr
		key  uint16
	}

	kvs, err := validatorsListDict.LoadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to load validators list dict: %w", err)
	}

	var totalWeight uint64
	var validatorsKeys = make([]validatorWithKey, len(kvs))
	for i, kv := range kvs {
		var val tlb.ValidatorAddr
		if err := tlb.LoadFromCell(&val, kv.Value); err != nil {
			return nil, fmt.Errorf("failed to parse validator addr: %w", err)
		}

		key, err := kv.Key.LoadUInt(16)
		if err != nil {
			return nil, fmt.Errorf("failed to parse validator key: %w", err)
		}

		totalWeight += val.Weight
		validatorsKeys[i].addr = &val
		validatorsKeys[i].key = uint16(key)
	}

	if definedWeight != nil && totalWeight != *definedWeight {
		return nil, fmt.Errorf("incorrect sum of weights")
	}

	if len(validatorsKeys) == 0 {
		return nil, fmt.Errorf("zero validators")
	}

	sort.Slice(validatorsKeys, func(i, j int) bool {
		return validatorsKeys[i].key < validatorsKeys[j].key
	})

	if validatorsNum > len(validatorsKeys) {
		validatorsNum = len(validatorsKeys)
	}

	var validators = make([]*tlb.ValidatorAddr, validatorsNum)

	for i := 0; i < validatorsNum; i++ {
		validators[i] = validatorsKeys[i].addr
	}

	return validators, nil
}
