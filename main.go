package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	force := false
	unindent := false
	flag.BoolVar(&force, "f", force, "force print even if invalid json")
	flag.BoolVar(&unindent, "u", unindent, "do not indent valid json")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "str2json - convert stringhified json to json\n")
		fmt.Fprintf(os.Stderr, "usage: %s [-f] [-u]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	buf, err := io.ReadAll(os.Stdin)
	check(err)

	buf = bytes.TrimSpace(buf)

	// Remove quotes.
	if bytes.HasPrefix(buf, []byte(`"`)) {
		buf = buf[1:]
	}
	if bytes.HasSuffix(buf, []byte(`"`)) {
		buf = buf[:len(buf)-1]
	}

	buf = bytes.ReplaceAll(buf, []byte(`\"`), []byte(`"`))

	if !json.Valid(buf) {
		if force {
			os.Stdout.Write(buf)
			return
		}
		fmt.Fprintf(os.Stderr, "invalid json, print anyway with -f (force)\n")
		os.Exit(1)
	}

	if unindent {
		os.Stdout.Write(buf)
		return
	}

	var out bytes.Buffer
	check(json.Indent(&out, buf, "", "  "))
	io.Copy(os.Stdout, &out)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
