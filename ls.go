package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
)

func getChildren(kv *api.KV, name string) ([]string, error) {
	name = strings.Trim(name, "/") + "/"
	if name == "/" {
		name = ""
	}
	pairs, _, err := kv.List(name, nil)
	if err != nil {
		return nil, err
	}
	childrenSet := make(map[string]struct{})
	for _, pair := range pairs {
		childPaths := strings.Split(strings.TrimPrefix(pair.Key, name), "/")
		if len(childPaths) == 0 {
			continue
		}
		childrenSet[childPaths[0]] = struct{}{}
	}
	children := []string{}
	for child := range childrenSet {
		children = append(children, child)
	}
	sort.Strings(children)
	return children, nil
}

func ls(kv *api.KV, name string, config *fsConfig) error {
	name = strings.Trim(name, "/")
	exist, err := consulFsFileExit(kv, name)
	if err != nil {
		return err
	}
	if exist {
		if !config.recurse {
			fmt.Println(name)
		}
		return nil
	}
	exist, err = consulFsFolderExist(kv, name)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("not found: %s", name)
	}
	children, err := getChildren(kv, name)
	if err != nil {
		return err
	}
	for _, child := range children {
		fmt.Println(name + "/" + child)
		if config.recurse {
			ls(kv, name+"/"+child, config)
		}
	}
	return nil
}

func cmdLs(kv *api.KV, args []string, config *fsConfig) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "USAGE: %s ls OPTIONS <remote>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "OPTIONS:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	err := ls(kv, args[0], config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
