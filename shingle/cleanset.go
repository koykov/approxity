package shingle

const (
	CleanSetPunct   = "!?.,;:'\"-—–()[]{}«»‹›"
	CleanSetSpecial = "@#$%^&*_+=\\|/~`<>"
	CleanSetAll     = CleanSetPunct + CleanSetSpecial
)
