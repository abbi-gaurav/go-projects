package word

import (
	"testing"
	"math/rand"
	"time"
	"unicode"
	"unicode/utf8"
)

func TestPalindrome(t *testing.T) {
	if !IsPalindrome("detartrated") {
		t.Error(`IsPalindrome("detartrated") = false`)
	}

	if !IsPalindrome("kayak") {
		t.Error(`IsPalindrome("kayak") = false`)
	}
}

func TestNonPalindrome(t *testing.T) {
	if IsPalindrome("palindrome") {
		t.Error(`IsPalindrome("palindrome") = true`)
	}
}

func TestFrenchPalindrome(t *testing.T) {
	if !IsPalindrome("été") {
		t.Error(`IsPalindrome(été) = false`)
	}
}

func TestCanalPalindrome(t *testing.T) {
	input := "A man, a plan, a canal : Panama"
	if !IsPalindrome(input) {
		t.Errorf(`IsPalindrome(%q) = false`, input)
	}
}

func TestIsPalindrome(t *testing.T) {
	var tests = [] struct {
		input    string
		expected bool
	}{
		{"", true},
		{"aa", true},
		{"aa", true},
		{"ab", false},
		{"kayak", true},
		{"detartrated", true},
		{"A man, a plan, a canal: Panama", true},
		{"Evil I did dwell; lewd did I live.", true},
		{"Able was I ere I saw Elba", true},
		{"été", true},
		{"Et se resservir, ivresse reste.", true},
		{"palindrome", false}, // non-palindrome
		{"desserts", false},   // semi-palindrome
	}

	for _, test := range tests {
		if actual := IsPalindrome(test.input); actual != test.expected {
			t.Errorf("IsPalindrome(%q) = %v expected [%v]", test.input, actual, test.expected)
		}
	}
}

func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		runes[i] = r
		runes[n-1-i] = r
	}

	return string(runes)
}

func TestRandomPalindromes(t *testing.T) {
	rng := getRNG(t)

	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}
func getRNG(t *testing.T) *rand.Rand {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed %d", seed)
	rng := rand.New(rand.NewSource(seed))
	return rng
}

func randomNonPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		for {
			c := rng.Intn(0x999)
			r := rune(c)
			r2 := rune(c + 1)
			if unicode.IsLetter(r) && unicode.IsLetter(r2) && unicode.ToLower(r) != unicode.ToLower(r2) {
				runes[i] = r
				runes[n-1-i] = r2
				break
			}
		}
	}
	return string(runes)
}

func TestRandomNonPalindrome(t *testing.T) {
	rng := getRNG(t)

	for i := 0; i < 1000; i++ {
		p := randomNonPalindrome(rng)
		if utf8.RuneCountInString(p) > 1 && IsPalindrome(p) {
			t.Errorf("IsPalindrome(%s) = true", p)
		}
	}
}
