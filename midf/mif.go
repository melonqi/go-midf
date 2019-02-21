package midf

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-spatial/geom"
)

/*
MifObj saves attributes and geo for one object
*/
type MifObj struct {
	Attributes []string
	Geo        geom.Geometry
}

/*
Mif saves info about mid/mif
*/
type Mif struct {
	Header  *MifHeader
	Objects []MifObj
}

/*
NewMif create new Mif.
*/
func NewMif() *Mif {
	return &Mif{Header: NewMifHeader()}
}

/*
HasColName use column name to query whether has column
*/
func (mif Mif) HasColName(name string) bool {
	lowerName := strings.ToLower(name)
	_, exists := mif.Header.NameMap[lowerName]
	return exists
}

/*
GetColPos get column index by name;
Because attributes saved as slice, use index will be fase.
return -1, if can't find col by name
*/
func (mif Mif) GetColPos(name string) int {
	if mif.HasColName(name) {
		return -1
	}

	lowerName := strings.ToLower(name)
	index, _ := mif.Header.NameMap[lowerName]
	return index
}

/*
AddColumn will add new column.
return -1, if existed; return 0, if success.
*/
func (mif *Mif) AddColumn(colName string, colType string) int {
	lowerName := strings.ToLower(colName)
	if mif.HasColName(colName) {
		return -1
	}

	mif.Header.ColNames = append(mif.Header.ColNames, colName)
	mif.Header.ColTypes = append(mif.Header.ColTypes, colType)

	mif.Header.NameMap[lowerName] = mif.Header.ColNum
	mif.Header.ColNum++
	for i := 0; i < len(mif.Objects); i++ {
		mif.Objects[i].Attributes = append(mif.Objects[i].Attributes, NullStr)
	}
	return 0
}

/*
Read will read mid/mif from file;
fileName: mid/mif name, without extension; For example, if you will read A.mid A.mif, just pass A as fileName.
midOnly: whether only have mid file
*/
func (mif *Mif) Read(fileName string, midOnly bool) int {
	mifFileName := fileName + ".mif"
	mifFile, err := os.Open(mifFileName)
	defer mifFile.Close()

	if err != nil {
		mifFileName = fileName + ".MIF"
		mifFile, err = os.Open(mifFileName)
		if err != nil {
			fmt.Printf("Open %s error, %s\n", mifFileName, err.Error())
			return -110
		}
	}
	//scanner have buffer size, this will cause imcomplete column
	scanner := bufio.NewScanner(mifFile)
	ret := mif.Header.getMifHeader(scanner)
	if ret < 0 {
		fmt.Println("get mif header failed")
	}
	ret = mif.getMifData(scanner)
	if err = scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return 0
}

func (mif *Mif) getMifData(scanner *bufio.Scanner) int {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		words := strings.Split(line, " ")
		geoType := strings.ToLower(words[0])
		switch {
		case geoType == "point":
			point, ret := getPointGeometry(words)
			if ret != 0 {
				return ret
			}
			mif.Objects = append(mif.Objects, MifObj{Geo: point})
		case geoType == "line":
			line, ret := getLineGeometry(words)
			if ret != 0 {
				return ret
			}
			mif.Objects = append(mif.Objects, MifObj{Geo: line})
		case geoType == "pline":
			multiLineString, ret := getMultiLineGeometry(words, scanner)
			if ret != 0 {
				return ret
			}
			mif.Objects = append(mif.Objects, MifObj{Geo: multiLineString})
		case geoType == "region":
			multiPolygon, ret := getRegionGeometry(words, scanner)
			if ret != 0 {
				return ret
			}
			mif.Objects = append(mif.Objects, MifObj{Geo: multiPolygon})
		case geoType == "pen" || strings.HasPrefix(geoType, "pen"):
			continue
		case geoType == "center":
			continue
		case geoType == "brush" || strings.HasPrefix(geoType, "brush"):
			continue
		case geoType == "symbol":
			continue
		case geoType == "rect":
			rect, ret := getRectGeometry(words)
			if ret != 0 {
				return ret
			}
			mif.Objects = append(mif.Objects, MifObj{Geo: rect})
		case geoType == "smooth":
			continue
		case geoType == "none":
			mif.Objects = append(mif.Objects, MifObj{Geo: nil})
		}
	}
	return 0
}

func getPointGeometry(words []string) (*geom.Point, int) {
	//POINT x y
	if len(words) < 3 {
		return nil, -1
	}
	var point geom.Point
	var err error
	point[0], err = strconv.ParseFloat(words[1], 64)
	if err != nil {
		return nil, -1
	}
	point[1], err = strconv.ParseFloat(words[2], 64)
	if err != nil {
		return nil, -1
	}
	return &point, 0
}

func getLineGeometry(words []string) (*geom.Line, int) {
	//LINE x1 y1 x2 y2
	if len(words) < 5 {
		return nil, -2
	}
	var line geom.Line
	var err error
	line[0][0], err = strconv.ParseFloat(words[1], 64)
	if err != nil {
		return nil, -2
	}
	line[0][1], err = strconv.ParseFloat(words[2], 64)
	if err != nil {
		return nil, -2
	}
	line[1][0], err = strconv.ParseFloat(words[3], 64)
	if err != nil {
		return nil, -2
	}
	line[1][1], err = strconv.ParseFloat(words[4], 64)
	if err != nil {
		return nil, -2
	}
	return &line, 0
}

func getMultiPoints(pointNum int, scanner *bufio.Scanner) [][2]float64 {
	line := make([][2]float64, pointNum)
	for j := 0; j < pointNum; j++ {
		scanner.Scan()
		coors := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		line[j][0], _ = strconv.ParseFloat(coors[0], 64)
		line[j][1], _ = strconv.ParseFloat(coors[1], 64)
	}
	return line
}

func getMultiLineGeometry(words []string, scanner *bufio.Scanner) (*geom.MultiLineString, int) {
	/*
	 * PLINE [ MULTIPLE numsections ]
	 numpts1
	 x1 y1
	 x2 y2
	 :
	 [ numpts2
	 x1 y1
	 x2 y2 ]
	 :
	*/

	var multiLineString geom.MultiLineString
	if len(words) == 3 {
		lineNum, err := strconv.Atoi(words[2])
		multiLines := make([][][2]float64, lineNum)
		if err != nil {
			return nil, -3
		}
		for i := 0; i < lineNum; i++ {
			scanner.Scan()
			pointNum, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				return nil, -3
			}
			line := getMultiPoints(pointNum, scanner)
			multiLines[i] = line
		}
		err = multiLineString.SetLineStrings(multiLines)
		if err != nil {
			return nil, -3
		}
	} else if len(words) == 2 {
		multiLines := make([][][2]float64, 1)
		pointNum, err := strconv.Atoi(words[1])
		if err != nil {
			return nil, -3
		}
		line := getMultiPoints(pointNum, scanner)
		multiLines[0] = line
		err = multiLineString.SetLineStrings(multiLines)
		if err != nil {
			return nil, -3
		}
	} else {
		return nil, -3
	}
	return &multiLineString, 0
}

func getRegionGeometry(words []string, scanner *bufio.Scanner) (*geom.MultiPolygon, int) {
	/*
	 REGION numpolygons
	 numpts1
	 x1 y1
	 x2 y2
	 :
	 [ numpts2
	 x1 y1
	 x2 y2 ]
	 :
	*/
	if len(words) != 2 {
		fmt.Println("wrong words size")
		return nil, -4
	}

	regionNum, err := strconv.Atoi(words[1])
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return nil, -4
	}

	var multiPolygon geom.MultiPolygon
	coors := make([][][][2]float64, regionNum)
	for i := 0; i < regionNum; i++ {
		scanner.Scan()
		pountNum, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return nil, -4
		}
		path := make([][][2]float64, 1)
		path[0] = getMultiPoints(pountNum, scanner)
		// if windingorder.OfPoints(path[0]...).IsCounterClockwise() {

		// }
		coors[i] = path
	}
	err = multiPolygon.SetPolygons(coors)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return nil, -4
	}

	return &multiPolygon, 0
}

func getRectGeometry(words []string) (*geom.Extent, int) {
	//Rect x0 y0 x1 y1
	if len(words) < 5 {
		return nil, -5
	}

	var err error
	var point geom.Point
	point[0], err = strconv.ParseFloat(words[1], 64)
	if err != nil {
		return nil, -5
	}
	point[1], err = strconv.ParseFloat(words[2], 64)
	if err != nil {
		return nil, -5
	}
	var rect geom.Extent
	rect.AddPoints(point)
	point[0], err = strconv.ParseFloat(words[3], 64)
	if err != nil {
		return nil, -5
	}
	point[1], err = strconv.ParseFloat(words[4], 64)
	if err != nil {
		return nil, -5
	}

	return &rect, 0
}
