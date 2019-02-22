package main

import (
	"fmt"
	"go-midf/midf"

	"github.com/go-spatial/geom"
)

func main() {
	mif := midf.NewMif()
	mif.Read("C_IndoorRegion", false)
	for _, obj := range mif.Objects {
		fmt.Println(obj.Attributes)
		switch obj.Geo.(type) {
		case *geom.Point:
			fmt.Println("*Point")
		case *geom.Line:
			fmt.Println("*Line")
		case *geom.MultiLineString:
			fmt.Println("*MultiLineString")
		case *geom.MultiPolygon:
			fmt.Println("*MultiPolygon")
		case *geom.Extent:
			fmt.Println("*Extent")
		}
	}
	// fmt.Println(midf.Split("1,2,3", ',', '"'))
	// fmt.Println(midf.Split("1,2,", ',', '"'))
	// fmt.Println(midf.Split("1,\"2,3\"", ',', '"'))
}
