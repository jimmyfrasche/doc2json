//Command doc2json(1) reads godoc(1) formatted text from stdin
//and writes json to stdout.
//
//Stdin is UTF-8 encoded plain text.
//The format of the text is exactly that as used in Go source documentation.
//However, the input is not expected to contain comment characters, see example.
//
//Stdout is minified JSON, in the format described below.
//
//JSON FORMAT
//
//The json is a list of objects.
//Each object has two fields: Kind and Lines.
//
//Kind is a string with the following possible values:
//	p   - paragraphs
//	h   - header
//	pre - preformatted text (likely code)
//
//Lines is a list of strings.
//
//EXAMPLE
//
//Note that doc2json does not pretty print its output, but it is
//done so here so as to make it easier to read.
//	$ cat <<EOF | doc2json
//	This is a paragraph.
//
//	This is a header
//
//	This is another paragraph.
//	With multiple lines.
//		This is preformatted
//		text
//	EOF
//	[
//		{
//			"Kind": "p",
//			"Lines": ["This is a paragraph.\n"]
//		},
//		{
//			"Kind": "h",
//			"Lines": ["This is a header"]
//		},
//		{
//			"Kind": "p",
//			"Lines": [
//				"This is another paragraph\n",
//				"With multiple lines.\n"
//			]
//		},
//		{
//			"Kind": "pre",
//			"Lines": [
//				"This is preformatted\n",
//				"text\n"
//			]
//		}
//	]
package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/jimmyfrasche/goutil"
)

type block struct {
	Kind  string
	Lines []string
}

func main() {
	log.SetFlags(0)

	if len(os.Args) != 1 {
		log.Fatalf("%s does not take any arguments", os.Args[0])
	}

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	doc := goutil.DocParse(string(stdin))

	blocks := make([]block, 0, len(doc))
	for _, d := range doc {
		var k string
		switch d.Kind {
		case goutil.Para:
			k = "p"
		case goutil.Head:
			k = "h"
		case goutil.Pre:
			k = "pre"
		default:
			log.Fatalln("godoc parser returned unrecognized output")
		}

		blocks = append(blocks, block{
			Kind:  k,
			Lines: d.Lines,
		})
	}

	bs, err := json.Marshal(blocks)
	if err != nil {
		log.Fatalln(err)
	}

	n, err := os.Stdout.Write(bs)
	if err != nil {
		log.Fatalln(err)
	}
	if err == nil && n < len(bs) {
		log.Fatalln(io.ErrShortWrite)
	}
}
