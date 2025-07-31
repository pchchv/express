package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Println("express [OPTIONS] <EXPRESSION> [DATA]")
	flag.PrintDefaults()
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func main() {

}
