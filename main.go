package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s OPTIONS ACTION ARGS\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "ACTION: get cat put ls delete")
	fmt.Fprintln(os.Stderr, "OPTIONS: ")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var (
		host    string
		scheme  string
		recurse bool
		force   bool
	)
	flag.StringVar(&host, "consul", "127.0.0.1:8500", "Consul API end point")
	flag.StringVar(&scheme, "scheme", "http", "Consul API scheme")
	flag.BoolVar(&recurse, "recurse", false, "work with put/get/delete/ls action")
	flag.BoolVar(&force, "force", false, "force put/get if folder exists")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		printUsage()
	}
	config := &fsConfig{
		recurse: recurse,
		force:   force,
	}

	client, err := api.NewClient(&api.Config{
		Address: host,
		Scheme:  scheme,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	kv := client.KV()

	switch args[0] {
	case "put":
		cmdPut(kv, args[1:], config)
	case "cat":
		cmdCat(kv, args[1:], config)
	case "get":
		cmdGet(kv, args[1:], config)
	case "ls":
		cmdLs(kv, args[1:], config)
	default:
		printUsage()
	}
}
