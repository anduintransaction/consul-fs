package main

import "fmt"

func getVersion() string {
	return "0.3"
}

func cmdVersion() {
	fmt.Println(getVersion())
}
