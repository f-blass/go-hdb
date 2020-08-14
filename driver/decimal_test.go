// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"math/big"
	"testing"
)

func testDecimalInfo(t *testing.T) {
	t.Logf("maximum decimal value %v", maxDecimal)
	t.Logf("~log2(10): %f", lg10)
}

func testDigits10(t *testing.T) {
	testData := []struct {
		x      *big.Int
		digits int
	}{
		{new(big.Int).SetInt64(0), 1},
		{new(big.Int).SetInt64(1), 1},
		{new(big.Int).SetInt64(9), 1},
		{new(big.Int).SetInt64(10), 2},
		{new(big.Int).SetInt64(99), 2},
		{new(big.Int).SetInt64(100), 3},
		{new(big.Int).SetInt64(999), 3},
		{new(big.Int).SetInt64(1000), 4},
		{new(big.Int).SetInt64(9999), 4},
		{new(big.Int).SetInt64(10000), 5},
		{new(big.Int).SetInt64(99999), 5},
		{new(big.Int).SetInt64(100000), 6},
		{new(big.Int).SetInt64(999999), 6},
		{new(big.Int).SetInt64(1000000), 7},
		{new(big.Int).SetInt64(9999999), 7},
		{new(big.Int).SetInt64(10000000), 8},
		{new(big.Int).SetInt64(99999999), 8},
		{new(big.Int).SetInt64(100000000), 9},
		{new(big.Int).SetInt64(999999999), 9},
		{new(big.Int).SetInt64(1000000000), 10},
		{new(big.Int).SetInt64(9999999999), 10},
	}

	for i, d := range testData {
		digits := digits10(d.x)
		if d.digits != digits {
			t.Fatalf("value %d int %s digits %d - expected digits %d", i, d.x, digits, d.digits)
		}
	}
}

func testConvertRat(t *testing.T) {
	testData := []struct {
		// in
		x      *big.Rat
		digits int
		minExp int
		maxExp int
		// out
		cmp *big.Int
		neg bool
		exp int
		df  decFlags
	}{
		{new(big.Rat).SetFrac64(0, 1), 3, -2, 2, new(big.Int).SetInt64(0), false, 0, 0},                              //convert 0
		{new(big.Rat).SetFrac64(1, 1), 3, -2, 2, new(big.Int).SetInt64(1), false, 0, 0},                              //convert 1
		{new(big.Rat).SetFrac64(1, 10), 3, -2, 2, new(big.Int).SetInt64(1), false, -1, 0},                            //convert 1/10
		{new(big.Rat).SetFrac64(1, 99), 3, -2, 2, new(big.Int).SetInt64(1), false, -2, dfNotExact},                   //convert 1/99
		{new(big.Rat).SetFrac64(1, 100), 3, -2, 2, new(big.Int).SetInt64(1), false, -2, 0},                           //convert 1/100
		{new(big.Rat).SetFrac64(1, 1000), 3, -2, 2, new(big.Int).SetInt64(1), false, -3, dfUnderflow},                //convert 1/1000
		{new(big.Rat).SetFrac64(10, 1), 3, -2, 2, new(big.Int).SetInt64(1), false, 1, 0},                             //convert 10
		{new(big.Rat).SetFrac64(100, 1), 3, -2, 2, new(big.Int).SetInt64(1), false, 2, 0},                            //convert 100
		{new(big.Rat).SetFrac64(1000, 1), 3, -2, 2, new(big.Int).SetInt64(10), false, 2, 0},                          //convert 1000
		{new(big.Rat).SetFrac64(10000, 1), 3, -2, 2, new(big.Int).SetInt64(100), false, 2, 0},                        //convert 10000
		{new(big.Rat).SetFrac64(100000, 1), 3, -2, 2, new(big.Int).SetInt64(100), false, 3, dfOverflow},              //convert 100000
		{new(big.Rat).SetFrac64(999999, 1), 3, -2, 2, new(big.Int).SetInt64(100), false, 4, dfNotExact | dfOverflow}, //convert 999999
		{new(big.Rat).SetFrac64(99999, 1), 3, -2, 2, new(big.Int).SetInt64(100), false, 3, dfNotExact | dfOverflow},  //convert 99999
		{new(big.Rat).SetFrac64(9999, 1), 3, -2, 2, new(big.Int).SetInt64(100), false, 2, dfNotExact},                //convert 9999
		{new(big.Rat).SetFrac64(99950, 1), 3, -2, 2, new(big.Int).SetInt64(100), false, 3, dfNotExact | dfOverflow},  //convert 99950
		{new(big.Rat).SetFrac64(99949, 1), 3, -2, 2, new(big.Int).SetInt64(999), false, 2, dfNotExact},               //convert 99949

		{new(big.Rat).SetFrac64(1, 3), 5, -5, 5, new(big.Int).SetInt64(33333), false, -5, dfNotExact}, //convert 1/3
		{new(big.Rat).SetFrac64(2, 3), 5, -5, 5, new(big.Int).SetInt64(66667), false, -5, dfNotExact}, //convert 2/3
		{new(big.Rat).SetFrac64(11, 2), 5, -5, 5, new(big.Int).SetInt64(55), false, -1, 0},            //convert 11/2

	}

	m := new(big.Int)

	for i := 0; i < 1; i++ { // use for performance tests
		for j, d := range testData {
			neg, exp, df := convertRatToDecimal(d.x, m, d.digits, d.minExp, d.maxExp)
			if m.Cmp(d.cmp) != 0 || neg != d.neg || exp != d.exp || df != d.df {
				t.Fatalf("converted %d value m %s neg %t exp %d df %b - expected m %s neg %t exp %d df %b", j, m, neg, exp, df, d.cmp, d.neg, d.exp, d.df)
			}
		}
	}
}

func TestDecimal(t *testing.T) {
	tests := []struct {
		name string
		fct  func(t *testing.T)
	}{
		{"decimalInfo", testDecimalInfo},
		{"digits10", testDigits10},
		{"convertRat", testConvertRat},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fct(t)
		})
	}
}
