package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
)

func cmdCat(kv *api.KV, args []string, config *fsConfig) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "USAGE: %s cat OPTIONS <file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "OPTIONS:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	file := strings.Trim(args[0], "/")
	p, _, err := kv.Get(file, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if p == nil {
		fmt.Fprintln(os.Stderr, "not found")
		os.Exit(1)
	}
	fmt.Print(string(p.Value))
}
