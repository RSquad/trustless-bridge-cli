package blockutils

import (
	"fmt"
	"math/big"

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

	switch t := set.Validators.(type) {
	case tlb.ValidatorSet:
		validatorsNum = int(t.Main)
		validatorsListDict = t.List
	case tlb.ValidatorSetExt:
		validatorsNum = int(t.Main)
		validatorsListDict = t.List
	default:
		return nil, 0, nil, fmt.Errorf("unknown validator set type")
	}

	type validatorWithKey struct {
		addr   *tlb.ValidatorAddr
		key    uint16
		weight uint64
	}

	kvs, err := validatorsListDict.LoadAll()
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to load validators list dict: %w", err)
	}

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

		validatorsKeys[i].addr = &val
		validatorsKeys[i].key = uint16(key)
		validatorsKeys[i].weight = val.Weight
	}

	if len(validatorsKeys) == 0 {
		return nil, 0, nil, fmt.Errorf("zero validators")
	}

	if validatorsNum > len(validatorsKeys) {
		validatorsNum = len(validatorsKeys)
	}

	var validators = make([]*tlb.ValidatorAddr, validatorsNum)

	var totalWeight uint64
	for i := 0; i < validatorsNum; i++ {
		validators[i] = validatorsKeys[i].addr
		totalWeight += validatorsKeys[i].weight
	}

	return validators, totalWeight, valDict.MustToCell().Hash(3), nil
}
