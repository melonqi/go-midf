package midf

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
CoordsysLL for longitude and latitude
*/
const CoordsysLL = "CoordSys Earth Projection 1, 0"

/*
CoordsysMC for mercator coordinate system
*/
const CoordsysMC = `CoordSys NonEarth Units "m" `

/*
NullStr for empty attribute
*/
const NullStr = `""`

/*
MifHeader save info from mif
like table column in sql
*/
type MifHeader struct {
	Version   int
	Charset   string
	Delimiter byte
	Coordsys  string
	ColNum    int
	Transform string
	ColNames  []string //to make sure every col has same struct, we use slice not map
	ColTypes  []string
	NameMap   map[string]int
}

/*
NewMifHeader create new mif header and init same values.
*/
func NewMifHeader() *MifHeader {
	return &MifHeader{
		Version:   300,
		Charset:   "WindowsSimpChinese",
		Delimiter: '\t',
		Coordsys:  CoordsysLL,
		NameMap:   make(map[string]int),
	}
}

/*
GetMifHeader gets mif header from scanner
*/
func (header *MifHeader) GetMifHeader(scanner *bufio.Scanner) int {
	colCnt := -2

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		words := strings.Split(line, " ")
		keyword := strings.ToLower(words[0])

		if colCnt < -1 {
			ret := header.getCommon(words, keyword, line, &colCnt)
			if ret != 0 {
				fmt.Printf("read common info error:%d\n", ret)
				return ret
			}
		} else {
			colCnt++
			if colCnt >= 0 && colCnt < header.ColNum {
				ret := header.getColumnInfo(words, keyword, line, &colCnt)
				if ret != 0 {
					fmt.Printf("read column info error:%d\n", ret)
					return ret
				}
			}
		}

		if keyword == "data" || keyword == "none" {
			break
		}
	}
	return len(header.NameMap)
}

func (header *MifHeader) getColumnInfo(words []string, keyword string, line string, colCnt *int) int {
	if len(words) < 2 {
		fmt.Printf("column type line has less than 2 words: %s\n", line)
		return -4
	}

	if keyword == "data" || keyword == "none" {
		fmt.Printf("invalid column type conf. need %d but find only %d\n", header.ColNum, colCnt)
		return -5
	}

	header.ColTypes[*colCnt] = words[1]
	lowerName := strings.ToLower(words[0])
	header.ColNames[*colCnt] = lowerName
	_, exists := header.NameMap[lowerName]
	if exists {
		fmt.Printf("re-define column: %s\n", lowerName)
		return -7
	}
	header.NameMap[lowerName] = *colCnt
	return 0
}

func (header *MifHeader) getCommon(words []string, keyword string, line string, colCnt *int) int {
	var err error
	switch keyword {
	case "version":
		if len(words) != 2 {
			fmt.Println("header version has not 2 word")
			return -1
		}
		header.Version, err = strconv.Atoi(words[1])
		if err != nil {
			fmt.Printf("parse header version error: %s\n", err.Error())
			return -1
		}
	case "charset":
		if len(words) != 2 {
			fmt.Println("header charset has not 2 word")
			return -2
		}
		header.Charset = strings.Trim(words[1], "\"")

	case "delimiter":
		if len(words) != 2 {
			fmt.Printf("header delimiter has not 2 word\n")
			return -6
		}
		header.Delimiter = words[1][1]
	//skip unique index
	// case "unique":
	// case "index":
	case "coordsys":
		header.Coordsys = line
	case "projection":
		header.Coordsys += " " + line
	case "transform":
		header.Transform = line
	case "columns":
		header.ColNum, err = strconv.Atoi(words[1])
		if err != nil {
			fmt.Printf("parse columns error: %s\n", err.Error())
			return -1
		}
		header.ColTypes = make([]string, header.ColNum)
		header.ColNames = make([]string, header.ColNum)
		//change status to parse column
		*colCnt = -1
	}

	return 0
}

/*
SetMifHeader will write mif header to file
*/
func (header *MifHeader) SetMifHeader(file *os.File) int {
	file.WriteString("Version " + strconv.Itoa(header.Version) + "\n")
	file.WriteString("Charset \"" + header.Charset + "\"\n")

	var delimiter []byte
	delimiter = append(delimiter, "Delimiter \""...)
	delimiter = append(delimiter, header.Delimiter)
	delimiter = append(delimiter, "\"\n"...)
	file.Write(delimiter)

	file.WriteString(header.Coordsys + "\n")
	if len(header.Transform) > 0 {
		file.WriteString(header.Transform + "\n")
	}
	if len(header.ColNames) != len(header.ColTypes) {
		return -2
	}
	file.WriteString("Columns " + strconv.Itoa(header.ColNum) + "\n")
	for i := 0; i < len(header.ColNames); i++ {
		file.WriteString("    " + header.ColNames[i] + " " + header.ColTypes[i] + "\n")
	}
	file.WriteString("Data\n")
	return 0
}
