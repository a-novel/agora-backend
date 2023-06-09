package security

import (
	"math/rand"
	"strings"
	"time"
)

var (
	// AlphabetLetters is an alphabet of latin non-accented characters.
	AlphabetLetters = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
)

// Random generates a random string of given size.
//
//	// Example output: iuhe
//	fmt.Println(security.Random(4))
//
// For more control over the set of characters, you may prefer RandomAlphabet.
// This method is a wrapper with AlphabetLetters as alphabet and an empty separator.
func Random(size int) string {
	return RandomAlphabet(AlphabetLetters, size, "")
}

// RandomAlphabet generates a random string from a given alphabet, and joins each element with the sep argument.
//
//	alphabet := []string{"foo", "bar", "qux"}
//
//	// Example output: qux-foo
//	fmt.Println(security.RandomAlphabet(alphabet, 2, "-"))
//
// For a set of random letters only, you may prefer Random.
func RandomAlphabet(alphabet []string, size int, sep string) string {
	rand.Seed(time.Now().UnixNano())

	alphabetSize := len(alphabet)
	output := make([]string, size)

	for i := 0; i < size; i++ {
		pick := rand.Intn(alphabetSize)
		output[i] = alphabet[pick]
	}

	return strings.Join(output, sep)
}
