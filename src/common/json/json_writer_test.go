package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapBracketsInFile(t *testing.T) {
	t.Run("one line", func(t *testing.T) {
		str := "{}[] not brackets ["
		expected := []Brackets{
			{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 1},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 2},
			{Type: CloseBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 3},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 18},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
	t.Run("one line, nested", func(t *testing.T) {
		str := "{bana: {1:1}, bana2:[1,2,3]}"
		expected := []Brackets{
			{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 7},
			{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 11},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 20},
			{Type: CloseBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 26},
			{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 27},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
	t.Run("multiple lines", func(t *testing.T) {
		str := "{\n}[] not \nbrackets \n["
		expected := []Brackets{
			{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: CloseBrackets, Shape: CurlyBrackets, Line: 2, CharIndex: 2},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 3},
			{Type: CloseBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 4},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 4, CharIndex: 21},
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
		expectedPairs := map[int]BracketPair{
			0: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 1}},
			2: {
				Open:  Brackets{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 2},
				Close: Brackets{Type: CloseBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 3}},
		}

		assert.Equal(t, expectedPairs, actualPairs)
	})

	t.Run("one line, nesting", func(t *testing.T) {
		str := "{bana: {1:1}, bana2:[1,2,3]}"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)

		expectedPairs := map[int]BracketPair{
			0: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 27}},
			7: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 7},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 11}},
			20: {
				Open:  Brackets{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 20},
				Close: Brackets{Type: CloseBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 26}},
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

		expectedPairs := map[int]BracketPair{
			0: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 2, CharIndex: 28}},
			7: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 7},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 11}},
			21: {
				Open:  Brackets{Type: OpenBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 21},
				Close: Brackets{Type: CloseBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 27}},
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
