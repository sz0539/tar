package main

import (
	"fmt"
	"github.com/sz0539/mytar/gtar"
)

func main() {

	t := &gtar.UnTar{
		Source: "",
		Target: "",
	}

	err := t.UnTarGzip()
	if err != nil {
		fmt.Println(err)
	}

}
