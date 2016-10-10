package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"
)

type fsConfig struct {
	recurse bool
	force   bool
}

func consulFsFileExit(kv *api.KV, name string) (bool, error) {
	name = strings.Trim(name, "/")
	if name == "" {
		return false, nil
	}
	p, _, err := kv.Get(name, nil)
	if err != nil {
		return false, err
	}
	if p != nil {
		return true, nil
	}
	return false, nil
}

func consulFsFolderExist(kv *api.KV, name string) (bool, error) {
	name = strings.Trim(name, "/") + "/"
	if name == "/" {
		name = ""
	}
	pairs, _, err := kv.List(name, nil)
	if err != nil {
		return false, err
	}
	if pairs != nil && len(pairs) > 0 {
		return true, nil
	}
	return false, nil
}

func consulListFolder(kv *api.KV, name string) ([]string, error) {
	name = strings.Trim(name, "/") + "/"
	if name == "/" {
		name = ""
	}
	pairs, _, err := kv.List(name, nil)
	if err != nil {
		return nil, err
	}
	if pairs == nil {
		return nil, fmt.Errorf("folder not exist: %s", name)
	}
	children := []string{}
	for _, pair := range pairs {
		list := strings.Split(strings.TrimPrefix(pair.Key, name), "/")
		if len(list) == 0 {
			continue
		}
		children = append(children, list[0])
	}
	return children, nil
}
