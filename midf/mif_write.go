package midf

import (
	"fmt"
	"os"
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

	return 0
}

func (mif *Mif) setMid(fileName string) int {
	return 0
}
