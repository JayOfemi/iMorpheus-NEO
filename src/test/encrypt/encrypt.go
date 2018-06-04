package encrypt

type Coder interface {
	Encode(input []byte) []byte

	Decode(input []byte) []byte
}