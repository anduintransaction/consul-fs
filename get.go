package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/consul/api"
)

func getFile(kv *api.KV, remote, local string, config *fsConfig) error {
	stat, err := os.Stat(local)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && stat.IsDir() {
		local = filepath.Join(local, filepath.Base(remote))
	}
	stat, err = os.Stat(local)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && stat.IsDir() {
		return fmt.Errorf("%s is a directory", local)
	}
	p, _, err := kv.Get(remote, nil)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("%s does not exist", remote)
	}
	return ioutil.WriteFile(local, p.Value, 0644)
}

func getFolder(kv *api.KV, remote, local string, config *fsConfig) error {
	if !config.recurse {
		return fmt.Errorf("%s is a directory, use -recurse to get", remote)
	}
	stat, err := os.Stat(local)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && !stat.IsDir() {
		return fmt.Errorf("%s is a normal file", local)
	} else if err == nil && stat.IsDir() {
		local = filepath.Join(local, filepath.Base(remote))
	}
	err = os.Mkdir(local, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		if !config.force {
			return fmt.Errorf("%s existed, use -force to force get", local)
		}
		err = os.RemoveAll(local)
		if err != nil {
			return err
		}
		err = os.Mkdir(local, 0755)
		if err != nil {
			return err
		}
	}
	children, err := consulListFolder(kv, remote)
	if err != nil {
		return err
	}
	for _, child := range children {
		err = get(kv, remote+"/"+child, filepath.Join(local, child), config)
		if err != nil {
			return err
		}
	}
	return nil
}

func get(kv *api.KV, remote, local string, config *fsConfig) error {
	fmt.Println(remote, local)
	exist, err := consulFsFolderExist(kv, remote)
	if err != nil {
		return err
	}
	if exist {
		return getFolder(kv, remote, local, config)
	}
	exist, err = consulFsFileExit(kv, remote)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("file %s does not exist", remote)
	}
	return getFile(kv, remote, local, config)
}

func cmdGet(kv *api.KV, args []string, config *fsConfig) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "USAGE: %s get OPTIONS <remote> <local>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "OPTIONS:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	remote := strings.Trim(args[0], "/")
	local := args[1]

	err := get(kv, remote, local, config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Success")
}
