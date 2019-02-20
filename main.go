package main

import (
	"fmt"
	"go-midf/midf"
)

func main() {
	mif := midf.NewMif()
	mif.Read("C_IndoorLine", false)
	fmt.Println(mif.Header.ColNames)
}
