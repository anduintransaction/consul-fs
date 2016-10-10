package main

import "fmt"

func getVersion() string {
	return "0.4"
}

func cmdVersion() {
	fmt.Println(getVersion())
}
