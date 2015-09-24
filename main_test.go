package main

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/trustmaster/go-aspell"
)

var (
	speller aspell.Speller
)

func TestCheckString(t *testing.T) {
	speller, err := aspell.NewSpeller(map[string]string{
		"lang": "en_US",
	})

	if err != nil {
		panic(err)
	}

	defer speller.Delete()

	cases := map[string][]result{
		"asdfasdf ":      {{0, 8, "asdfasdf"}},
		"asdfasdf":       {{0, 8, "asdfasdf"}},
		"hello asdfasdf": {{6, 14, "asdfasdf"}},
		" asdfasdf ":     {{1, 9, "asdfasdf"}},
		"AsdfasdfHello":  {{0, 8, "Asdfasdf"}},
	}

	splitCamelCase = true

	for input, output := range cases {
		realOutput := checkString(speller, input)

		if len(realOutput) != len(output) {
			t.Errorf("Poop %v %v %v", input, output, realOutput)
		} else {
			for i := 0; i < len(output); i++ {
				if !realOutput[i].SameAs(output[i]) {
					t.Errorf("Poop %v %v %v", input, output, realOutput)
				}
			}
		}
	}
}

func ExampleCheckComment() {
	splitCamelCase = false
	speller, err := aspell.NewSpeller(map[string]string{
		"lang": "en_US",
	})

	if err != nil {
		panic(err)
	}

	comment := &ast.Comment{Slash: 1, Text: "// Spellign Msitake spellingMistake SpellingMistake"}
	fs := token.NewFileSet()
	fs.AddFile("/foo/bar/baz.go", 1, 2048)

	checkComment(speller, fs, comment)
	// Output:
	///foo/bar/baz.go:1:4: Spellign
	///foo/bar/baz.go:1:13: Msitake
	///foo/bar/baz.go:1:21: spellingMistake
	///foo/bar/baz.go:1:37: SpellingMistake

}

func ExampleCheckComment_camel() {
	splitCamelCase = true
	speller, err := aspell.NewSpeller(map[string]string{
		"lang": "en_US",
	})

	if err != nil {
		panic(err)
	}

	comment := &ast.Comment{Slash: 1, Text: "// Spellign Msitake spellingMistake SpellingMistake"}
	fs := token.NewFileSet()
	fs.AddFile("/foo/bar/baz.go", 1, 2048)

	checkComment(speller, fs, comment)
	// Output:
	///foo/bar/baz.go:1:4: Spellign
	///foo/bar/baz.go:1:13: Msitake
}
