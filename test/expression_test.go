package test

import (
	"testing"
)

func TestArithmeticExpressions(t *testing.T) {
	testExpressions(t, []testCase{
		{"5 + 5", "10"},
		{"5 - 3", "2"},
		{"-3", "-3"},
		{"4 * 2", "8"},
		{"8 / 2", "4"},
		{"5 + 2 * 3", "11"},
		{"(5 + 2) * 3", "21"},
		{"5 + 30 / 3", "15"},
		{"(5 + 3) / 2", "4"},
		{"100 % 19", "5"},
	})
}

func TestBooleanExpressions(t *testing.T) {
	testExpressions(t, []testCase{
		{"5 > 5", "false"},
		{"5 >= 5", "true"},
		{"7 > 5", "true"},
		{"11 < 11", "false"},
		{"11 <= 11", "true"},
		{"11 <= 11", "true"},
		{"13 != 11", "true"},
		{"13 != 13", "false"},
		{"17 == 13", "false"},
		{"17 == 17", "true"},
		{"nil == nil", "true"},
		{"nil or 19", "19"},
		{"nil and 19", "nil"},
		{"23 or nil", "23"},
		{"23 and nil", "nil"},
		{"true ? 1 : 0", "1"},
		{"29 > 31 ? false : 31 > 29 ? true : false", "true"},
	})
}
