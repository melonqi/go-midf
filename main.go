package main

import (
	"fmt"
	"go-midf/midf"

	"github.com/go-spatial/geom"
)

func main() {
	mif := midf.NewMif()
	mif.Read("test")
	idPos := mif.GetColPos("id")
	for _, obj := range mif.Objects {
		fmt.Println(obj.Attributes[idPos])
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
	// mif.AddColumn("id", "Char(16)")
	// mif.Header.Coordsys = midf.CoordsysLL
	// var obj midf.MifObj
	// obj.Attributes = append(obj.Attributes, "123456")
	// var geo geom.Point
	// geo[0] = 116.0
	// geo[1] = 40.0
	// obj.Geo = &geo
	// mif.Objects = append(mif.Objects, obj)
	// mif.Write("test")
	// mif.Read("C_IndoorRegion")
	// mif.Write("test")

	/*
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
	*/
	// fmt.Println(midf.Split("1,2,3", ',', '"'))
	// fmt.Println(midf.Split("1,2,", ',', '"'))
	// fmt.Println(midf.Split("1,\"2,3\"", ',', '"'))
}
