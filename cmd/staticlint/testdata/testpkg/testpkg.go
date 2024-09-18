package main

import "os"

func callOsExitFunc() {
	os.Exit(0)
}

func main() {
	os.Exit(0) // want "directly call os.Exit"
}
