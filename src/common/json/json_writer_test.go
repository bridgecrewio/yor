package json

import (
	"github.com/bridgecrewio/yor/src/common/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapBracketsInFile(t *testing.T) {
	t.Run("one line", func(t *testing.T) {
		str := "{}[] not brackets ["
		expected := []types.Brackets{
			{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 1},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 2},
			{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 3},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 18},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
	t.Run("one line, nested", func(t *testing.T) {
		str := "{bana: {1:1}, bana2:[1,2,3]}"
		expected := []types.Brackets{
			{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 7},
			{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 11},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 20},
			{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 26},
			{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 27},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
	t.Run("multiple lines", func(t *testing.T) {
		str := "{\n}[] not \nbrackets \n["
		expected := []types.Brackets{
			{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 2, CharIndex: 2},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 2, CharIndex: 3},
			{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 2, CharIndex: 4},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 4, CharIndex: 21},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
}

func TestGetBracketsPairs(t *testing.T) {
	t.Run("one line, no nesting", func(t *testing.T) {
		str := "{}[] not brackets"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)
		expectedPairs := map[int]types.BracketPair{
			0: {
				Open:  types.Brackets{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
				Close: types.Brackets{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 1}},
			2: {
				Open:  types.Brackets{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 2},
				Close: types.Brackets{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 3}},
		}

		assert.Equal(t, expectedPairs, actualPairs)
	})

	t.Run("one line, nesting", func(t *testing.T) {
		str := "{bana: {1:1}, bana2:[1,2,3]}"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)

		expectedPairs := map[int]types.BracketPair{
			0: {
				Open:  types.Brackets{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
				Close: types.Brackets{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 27}},
			7: {
				Open:  types.Brackets{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 7},
				Close: types.Brackets{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 11}},
			20: {
				Open:  types.Brackets{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 20},
				Close: types.Brackets{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 26}},
		}
		for index, pair := range expectedPairs {
			actualPair, ok := actualPairs[index]
			if !ok {
				t.Errorf("expected to get pair in index %d", index)
			}
			assert.Equal(t, pair, actualPair)
		}
	})
	t.Run("multiple lines with nesting", func(t *testing.T) {
		str := "{bana: {1:1},\n bana2:[1,2,3]}"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)

		expectedPairs := map[int]types.BracketPair{
			0: {
				Open:  types.Brackets{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
				Close: types.Brackets{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 2, CharIndex: 28}},
			7: {
				Open:  types.Brackets{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 7},
				Close: types.Brackets{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 11}},
			21: {
				Open:  types.Brackets{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 2, CharIndex: 21},
				Close: types.Brackets{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 2, CharIndex: 27}},
		}
		for index, pair := range expectedPairs {
			actualPair, ok := actualPairs[index]
			if !ok {
				t.Errorf("expected to get pair in index %d", index)
			}
			assert.Equal(t, pair, actualPair)
		}
	})
}
