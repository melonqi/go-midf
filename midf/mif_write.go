package midf

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-spatial/geom"
)

func (mif *Mif) Write(fileName string) bool {
	ret := mif.setMif(fileName)
	if ret < 0 {
		fmt.Printf("Write %s mif failed\n", fileName)
		return false
	}
	ret = mif.setMid(fileName)
	if ret < 0 {
		fmt.Printf("Write %s mid failed\n", fileName)
		return false
	}
	return true
}

func (mif *Mif) setMif(fileName string) int {
	mifFileName := fileName + ".mif"
	mifFile, err := os.Create(mifFileName)
	defer mifFile.Close()
	if err != nil {
		fmt.Printf("Open %s error, %s\n", mifFileName, err.Error())
		return -1
	}

	ret := mif.Header.SetMifHeader(mifFile)
	if ret < 0 {
		fmt.Println("set mif header failed")
		return -1
	}
	mif.setMifData(mifFile)
	return 0
}

func (mif *Mif) setMifData(file *os.File) {
	for i := 0; i < len(mif.Objects); i++ {
		if mif.Objects[i].Geo == nil {
			file.WriteString("None\n")
		} else {
			switch mif.Objects[i].Geo.(type) {
			case *geom.Point:
				pt := mif.Objects[i].Geo.(*geom.Point)
				file.WriteString("Point " + strconv.FormatFloat(pt[0], 'g', 12, 64))
				file.WriteString(" " + strconv.FormatFloat(pt[1], 'g', 12, 64) + "\n")
				file.WriteString("\n")
			case *geom.MultiLineString:
				multiLineString := mif.Objects[i].Geo.(*geom.MultiLineString).LineStrings()
				partSize := len(multiLineString)
				file.WriteString("Pline MULTIPLE " + strconv.Itoa(partSize) + "\n")
				for j := 0; j < partSize; j++ {
					ptSize := len(multiLineString[j])
					file.WriteString("    " + strconv.Itoa(ptSize) + "\n")
					for k := 0; k < ptSize; k++ {
						file.WriteString(strconv.FormatFloat(multiLineString[j][k][0], 'g', 12, 64))
						file.WriteString(" ")
						file.WriteString(strconv.FormatFloat(multiLineString[j][k][1], 'g', 12, 64))
						file.WriteString("\n")
					}
				}
				file.WriteString("    Pen (1,2,0)\n")
			case *geom.MultiPolygon:
				multiPolygon := mif.Objects[i].Geo.(*geom.MultiPolygon).Polygons()
				partSize := len(multiPolygon)
				file.WriteString("Region " + strconv.Itoa(partSize) + "\n")
				for j := 0; j < partSize; j++ {
					ptSize := len(multiPolygon[j][0])
					file.WriteString("  " + strconv.Itoa(ptSize) + "\n")
					for k := 0; k < ptSize; k++ {
						file.WriteString(strconv.FormatFloat(multiPolygon[j][0][k][0], 'g', 12, 64))
						file.WriteString(" ")
						file.WriteString(strconv.FormatFloat(multiPolygon[j][0][k][1], 'g', 12, 64))
						file.WriteString("\n")
					}
				}
			}
		}
	}
}

func (mif *Mif) setMid(fileName string) int {
	midFileName := fileName + ".mid"
	midFile, err := os.Create(midFileName)
	defer midFile.Close()
	if err != nil {
		fmt.Printf("Open %s error, %s\n", midFileName, err.Error())
		return -1
	}

	var delimiter []byte
	delimiter = append(delimiter, mif.Header.Delimiter)
	sep := string(delimiter)
	for i := 0; i < len(mif.Objects); i++ {
		midFile.WriteString(strings.Join(mif.Objects[i].Attributes, sep))
		if i != len(mif.Objects)-1 {
			midFile.WriteString("\n")
		}
	}
	return 0
}
