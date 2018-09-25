package cipher

import "unicode"

type caesarCipher struct {}

func NewCaesarCipher() (*caesarCipher, error) {
	return &caesarCipher{}, nil
}

func (c *caesarCipher) Run() error {
	return nil
}

func Encode(s string, k int) string {
	arr := make([]byte, len(s), len(s))
	key := uint8(k % 26)
	var offset uint8
	for i := 0; i < len(s); i++ {
		c := s[i]
		if unicode.IsLetter(rune(c)) {
			if unicode.IsUpper(rune(c)) {
				offset = 65
			} else {
				offset = 97
			}
			n := c - offset
			n = (n + key) % 26
			arr[i] = n + offset
		} else {
			arr[i] = c
		}
	}
	str := string(arr)
	return str
}