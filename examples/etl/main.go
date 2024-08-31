package main

import "github.com/reugn/tpack"

func main() {
	tpack.NewPackerStd(tpack.NewProcessor(
		doETL,
	)).Execute()
}
