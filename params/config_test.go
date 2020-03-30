// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"
	"reflect"
	"testing"
)

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new *ChainConfig
		head        uint64
		wantErr     *ConfigCompatError
	}
	var maxCodeSizev1, maxCodeSizev2 []MaxCodeSizeStruct
	maxCodeSize0 := MaxCodeSizeStruct{big.NewInt(0), 24}
	maxCodeSize10 := MaxCodeSizeStruct{big.NewInt(10), 32}
	maxCodeSize20 := MaxCodeSizeStruct{big.NewInt(20), 48}

	maxCodeSizev1 = append(maxCodeSizev1, maxCodeSize0)
	maxCodeSizev1 = append(maxCodeSizev1, maxCodeSize10)

	maxCodeSizev2 = append(maxCodeSizev2, maxCodeSize0)
	maxCodeSizev2 = append(maxCodeSizev2, maxCodeSize20)

	tests := []test{
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 0, wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 100, wantErr: nil},
		{
			stored:  &ChainConfig{EIP150Block: big.NewInt(10)},
			new:     &ChainConfig{EIP150Block: big.NewInt(20)},
			head:    9,
			wantErr: nil,
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &ChainConfig{HomesteadBlock: nil},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    nil,
				RewindTo:     0,
			},
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &ChainConfig{HomesteadBlock: big.NewInt(1)},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    big.NewInt(1),
				RewindTo:     0,
			},
		},
		{
			stored: &ChainConfig{HomesteadBlock: big.NewInt(30), EIP150Block: big.NewInt(10)},
			new:    &ChainConfig{HomesteadBlock: big.NewInt(25), EIP150Block: big.NewInt(20)},
			head:   25,
			wantErr: &ConfigCompatError{
				What:         "EIP150 fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
		{
			stored:  &ChainConfig{Istanbul: &IstanbulConfig{Ceil2Nby3Block: big.NewInt(10)}},
			new:     &ChainConfig{Istanbul: &IstanbulConfig{Ceil2Nby3Block: big.NewInt(20)}},
			head:    4,
			wantErr: nil,
		},
		{
			stored: &ChainConfig{Istanbul: &IstanbulConfig{Ceil2Nby3Block: big.NewInt(10)}},
			new:    &ChainConfig{Istanbul: &IstanbulConfig{Ceil2Nby3Block: big.NewInt(20)}},
			head:   30,
			wantErr: &ConfigCompatError{
				What:         "Ceil 2N/3 fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
		{
			stored: &ChainConfig{MaxCodeSize : maxCodeSizev1},
			new:    &ChainConfig{MaxCodeSize: maxCodeSizev2},
			head:   30,
			wantErr: &ConfigCompatError{
				What:         "max code size change fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
		//{
		//	stored:  &ChainConfig{MaxCodeSizeChangeBlock:big.NewInt(10)},
		//	new:     &ChainConfig{MaxCodeSizeChangeBlock:big.NewInt(20)},
		//	head:    4,
		//	wantErr: nil,
		//},
		{
			stored: &ChainConfig{QIP714Block:big.NewInt(10)},
			new:    &ChainConfig{QIP714Block:big.NewInt(20)},
			head:   30,
			wantErr: &ConfigCompatError{
				What:         "permissions fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
		{
			stored:  &ChainConfig{QIP714Block:big.NewInt(10)},
			new:     &ChainConfig{QIP714Block:big.NewInt(20)},
			head:    4,
			wantErr: nil,
		},

	}

	for _, test := range tests {
		err := test.stored.CheckCompatible(test.new, test.head, false)
		if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		}
	}
}
