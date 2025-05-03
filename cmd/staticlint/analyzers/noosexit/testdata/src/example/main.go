package main

import "os"

func main() {
	os.Exit(1) // want "os.Exit called in main package"
}

func other() {
	os.Exit(1)
}
