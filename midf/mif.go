package midf

import (
	"bufio"
	"fmt"
	"os"
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
	if err = scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return 0
}
