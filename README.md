# go-midf
MapInfo mid/mif reader and writer go library

test.mid/mif have one attribute which name is 'id';

in this mid/mif,  there is one obj, which id=123456, location is (116.0,40.0)

## test.mif

```
Delimiter "	"
CoordSys Earth Projection 1, 0
Columns 1
    id Char(16)
Data
Point 116 40
```

## test.mid 

```
123456
```



## Write mid/mif

```go
mif := midf.NewMif()
mif.AddColumn("id", "Char(16)")
mif.Header.Coordsys = midf.CoordsysLL
var obj midf.MifObj
obj.Attributes = append(obj.Attributes, "123456")
var geo geom.Point
geo[0] = 116.0
geo[1] = 40.0
obj.Geo = &geo
mif.Objects = append(mif.Objects, obj)
mif.Write("test")
```

## Read midf

```go
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
```