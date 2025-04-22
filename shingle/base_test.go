package shingle

type stage struct {
	name    string
	text    string
	tokens  map[uint][]string
	ctokens map[uint][]string
}
