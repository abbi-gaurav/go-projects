package word

import "unicode"

func IsPalindrome(string string) bool {
	var letters []rune
	for _, r := range string {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	for i := range letters {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}
