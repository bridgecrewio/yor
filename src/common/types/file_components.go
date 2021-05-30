package types

type BracketShape int
type BracketDirection int

const (
	CurlyBrackets BracketShape = iota + 1
	SquareBrackets

	OpenBrackets BracketDirection = iota + 1
	CloseBrackets
)

type Brackets struct {
	Type      BracketDirection
	Shape     BracketShape
	Line      int
	CharIndex int
}

type BracketPair struct {
	Open  Brackets
	Close Brackets
}
