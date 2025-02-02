package blockutils

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/rsquad/trustless-bridge-cli/internal/tonclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func ExtractMainValidators(block *tlb.Block, tonClient *tonclient.TonClient) ([]*tlb.ValidatorAddr, uint64, []byte, error) {
	c, err := block.Extra.Custom.ConfigParams.Config.Params.LoadValueByIntKey(big.NewInt(34))
	if err != nil {
		return nil, 0, nil, err
	}
	valDict := c.MustLoadRef()
	var set tlb.ValidatorSetAny
	if err = tlb.LoadFromCell(&set, valDict); err != nil {
		return nil, 0, nil, err
	}

	var validatorsNum int
	var validatorsListDict *cell.Dictionary

	var definedWeight *uint64
	switch t := set.Validators.(type) {
	case tlb.ValidatorSet:
		validatorsNum = int(t.Main)
		validatorsListDict = t.List
	case tlb.ValidatorSetExt:
		definedWeight = &t.TotalWeight
		validatorsNum = int(t.Main)
		validatorsListDict = t.List
	default:
		return nil, 0, nil, fmt.Errorf("unknown validator set type")
	}

	type validatorWithKey struct {
		addr *tlb.ValidatorAddr
		key  uint16
	}

	kvs, err := validatorsListDict.LoadAll()
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to load validators list dict: %w", err)
	}

	var totalWeight uint64
	var validatorsKeys = make([]validatorWithKey, len(kvs))
	for i, kv := range kvs {
		var val tlb.ValidatorAddr
		if err := tlb.LoadFromCell(&val, kv.Value); err != nil {
			return nil, 0, nil, fmt.Errorf("failed to parse validator addr: %w", err)
		}

		key, err := kv.Key.LoadUInt(16)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("failed to parse validator key: %w", err)
		}

		totalWeight += val.Weight
		validatorsKeys[i].addr = &val
		validatorsKeys[i].key = uint16(key)
	}

	if definedWeight != nil && totalWeight != *definedWeight {
		return nil, 0, nil, fmt.Errorf("incorrect sum of weights")
	}

	if len(validatorsKeys) == 0 {
		return nil, 0, nil, fmt.Errorf("zero validators")
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

	return validators, totalWeight, valDict.MustToCell().Hash(3), nil
}
