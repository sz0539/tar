package main

import (
	"fmt"
	"github.com/sz0539/tar/untar"
)

func main() {

	t := &untar.UnTar{
		Source: "",
		Target: "",
	}

	err := t.UnTarGzip()
	if err != nil {
		fmt.Println(err)
	}

}
