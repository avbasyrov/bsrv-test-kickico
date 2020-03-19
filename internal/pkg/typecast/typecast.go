package typecast

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"reflect"
	"strings"
)

func FromString(abiType abi.Type, value string) (interface{}, error) {
	switch abiType.T {
	case abi.BoolTy:
		value = strings.ToLower(value)
		if value == "1" || value == "true" {
			return true, nil
		} else if value == "0" || value == "false" {
			return false, nil
		} else {
			return nil, errors.New(fmt.Sprintf("expected bool value, but got: %s", value))
		}
	case abi.AddressTy:
		if !common.IsHexAddress(value) {
			return nil, errors.New(fmt.Sprintf("invalid ethereum address: %s", value))
		}
		if has0xPrefix(value) {
			value = value[2:]
		}
		decoded, err := hex.DecodeString(value)
		if err != nil {
			return nil, err
		}
		if len(decoded) != 20 {
			return nil, errors.New("WTF")
		}
		var address [20]byte
		for i, v := range decoded {
			address[i] = v
		}

		return address, nil
	case abi.UintTy:
		switch abiType.Kind {
		case reflect.Ptr: // uint256 & int256
			bigInt, success := big.NewInt(0).SetString(value, 10)
			if !success {
				return nil, errors.New(fmt.Sprintf("ivalid uint256 number: '%s'", value))
			}
			return bigInt, nil
		default:
			return nil, errors.New(fmt.Sprintf("unsupported integer type '%s'", abiType.String()))
		}
	case abi.BytesTy:
		return []byte(value), nil
	}

	return nil, errors.New(fmt.Sprintf("unsupported type '%s'", abiType.String()))
}

func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}
