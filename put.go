package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/consul/api"
)

func putFile(kv *api.KV, reader io.Reader, target string) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	pair := &api.KVPair{
		Key:   target,
		Value: content,
	}
	_, err = kv.Put(pair, nil)
	return err
}

func putFolder(kv *api.KV, config *fsConfig, local, remote string) error {
	basename := filepath.Base(local)
	if !config.recurse {
		return fmt.Errorf("%s is a directory, use -recurse to put", local)
	}
	exist, err := consulFsFolderExist(kv, remote)
	if err != nil {
		return err
	}
	if exist {
		// put inside
		remote += "/" + basename
	}
	exist, err = consulFsFileExit(kv, remote)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("%s is a normal file", remote)
	}
	exist, err = consulFsFolderExist(kv, remote)
	if err != nil {
		return err
	}
	if exist {
		if !config.force {
			return fmt.Errorf("%s existed, use -force to force put", remote)
		}
		_, err = kv.DeleteTree(remote+"/", nil)
		if err != nil {
			return err
		}
	}
	localAbsPath, err := filepath.Abs(local)
	if err != nil {
		return err
	}
	err = filepath.Walk(local, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(path, localAbsPath+"/")
		reader, err := os.Open(path)
		if err != nil {
			return err
		}
		defer reader.Close()
		return putFile(kv, reader, remote+"/"+relPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func cmdPut(kv *api.KV, args []string, config *fsConfig) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "USAGE: %s put OPTIONS <remote> (read content from stdin)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s put OPTIONS <local file or folder> <remote>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "OPTIONS:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if len(args) >= 2 {
		local := args[0]
		remote := strings.Trim(args[1], "/")
		localStat, err := os.Stat(local)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if localStat.IsDir() {
			err = putFolder(kv, config, local, remote)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else {
			exist, err := consulFsFolderExist(kv, remote)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if exist {
				// put inside
				remote += "/" + filepath.Base(local)
			}
			reader, err := os.Open(local)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			err = putFile(kv, reader, remote)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	} else {
		remote := strings.Trim(args[0], "/")
		exist, err := consulFsFolderExist(kv, remote)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if exist {
			fmt.Fprintln(os.Stderr, "remote folder existed")
			os.Exit(1)
		}
		err = putFile(kv, os.Stdin, remote)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("Success")
}
