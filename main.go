package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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

type ConsulConfig struct {
	Scheme   string `json:"scheme"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func readConfig() (*ConsulConfig, error) {
	b, err := ioutil.ReadFile(os.Getenv("HOME") + "/.consul/config.json")
	if err != nil {
		return &ConsulConfig{
			Scheme: "http",
			Host:   "127.0.0.1:8500",
		}, nil
	}
	var config ConsulConfig
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	fileConfig, err := readConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var (
		host     string
		scheme   string
		user     string
		password string
		recurse  bool
		force    bool
	)
	flag.StringVar(&host, "consul", fileConfig.Host, "Consul API end point")
	flag.StringVar(&scheme, "scheme", fileConfig.Scheme, "Consul API scheme")
	flag.StringVar(&user, "user", fileConfig.User, "Consul API scheme")
	flag.StringVar(&password, "password", fileConfig.Password, "Consul API scheme")
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

	consulConfig := &api.Config{
		Address: host,
		Scheme:  scheme,
	}
	if user != "" && password != "" {
		consulConfig.HttpAuth = &api.HttpBasicAuth{
			Username: user,
			Password: password,
		}
	}

	client, err := api.NewClient(consulConfig)
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
	case "version":
		cmdVersion()
	default:
		printUsage()
	}
}
