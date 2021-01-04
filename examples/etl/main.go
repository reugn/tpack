package main

import "github.com/reugn/tpack"

func main() {
	tpack.NewPackerStd(tpack.NewFunctionProcessor(
		doETL,
	)).Execute()
}
