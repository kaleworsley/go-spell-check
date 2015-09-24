package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"unicode"

	"github.com/trustmaster/go-aspell"
)

var (
	splitCamelCase bool
	lang           string
)

type result struct {
	start int
	end   int
	str   string
}

func (r result) SameAs(o result) bool {
	return r.start == o.start && r.end == o.end && r.str == o.str
}

func init() {
	flag.BoolVar(&splitCamelCase, "camel", false, "treak CamelCased words as seperate words")
	flag.StringVar(&lang, "lang", "en_US", "aspell language")
}

func checkString(speller aspell.Speller, str string) (output []result) {
	var (
		inWord        = false
		wordStart     int
		wordEnd       int
		prevLowerCase bool
	)

	for i, r := range str {
		if unicode.IsLetter(r) || r == '\'' {
			if !inWord {
				wordStart = len(str[:i])
				inWord = true
			} else {
				if splitCamelCase {
					if unicode.IsLetter(r) {
						if unicode.IsLower(r) {
							prevLowerCase = true
						} else {
							if prevLowerCase {
								wordEnd = len(str[:i])

								word := str[wordStart:wordEnd]
								if !speller.Check(word) {
									output = append(output, result{wordStart, wordEnd, word})
								}
								inWord = true
								wordStart = i
								prevLowerCase = false
							}
						}
					}
				}
			}
		} else {
			if inWord {
				wordEnd = len(str[:i])

				word := str[wordStart:wordEnd]
				if !speller.Check(word) {
					output = append(output, result{wordStart, wordEnd, word})
				}
				inWord = false
				prevLowerCase = false
			}
		}

	}
	if inWord {
		wordEnd = len(str)

		word := str[wordStart:wordEnd]
		if !speller.Check(word) {
			output = append(output, result{wordStart, wordEnd, word})
		}
		inWord = false
		prevLowerCase = false
	}

	return output
}

func checkComment(speller aspell.Speller, fs *token.FileSet, comment *ast.Comment) {
	results := checkString(speller, comment.Text)

	for _, res := range results {
		pos := fs.Position(comment.Slash + token.Pos(res.start))
		fmt.Printf("%v:%v:%v: %v\n", pos.Filename, pos.Line, pos.Column, res.str)
	}
}

func main() {
	flag.Parse()

	var dir string
	if len(flag.Args()) == 2 {
		dir = flag.Args()[1]
	} else {
		dir, _ = os.Getwd()
	}

	speller, err := aspell.NewSpeller(map[string]string{
		"lang": lang,
	})

	if err != nil {
		panic(err)
	}

	defer speller.Delete()

	fs := token.NewFileSet()

	pkgs, err := parser.ParseDir(fs, dir, nil, parser.ParseComments)

	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			for _, cg := range f.Comments {
				for _, comment := range cg.List {
					checkComment(speller, fs, comment)
				}
			}
		}
	}
}
