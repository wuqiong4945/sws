package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
)

var cfg *ini.File
var columnSettings []*ini.Key
var foFolder string = "fo"

func main() {
	var err error
	cfg, err = ini.Load("sws.ini")
	printError(err)
	srcFolder := cfg.Section("general").Key("srcfolder").MustString("src")
	swsFolder := cfg.Section("general").Key("swsfolder").MustString("sws")

	os.MkdirAll(foFolder, os.ModeDir|os.ModePerm)
	createSws(srcFolder, swsFolder)

	os.RemoveAll(foFolder)
}

func createSws(srcFolder, swsFolder string) {
	os.MkdirAll(swsFolder, os.ModeDir|os.ModePerm)

	fmt.Printf("\n-dir %+v\n", srcFolder)
	srcFileInfoList, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		printError(err)
		return
	}

	// main loop, deal with all src xml files
	for _, srcFileInfo := range srcFileInfoList {
		if srcFileInfo.IsDir() {
			if !strings.HasPrefix(srcFileInfo.Name(), ".") {
				createSws(srcFolder+"/"+srcFileInfo.Name(), swsFolder+"/"+srcFileInfo.Name())
			}
			continue
		}
		if !strings.HasSuffix(srcFileInfo.Name(), ".xml") {
			continue
		}
		// fmt.Println(srcFileInfo.Name())
		fileName := strings.TrimSuffix(srcFileInfo.Name(), ".xml")

		srcFile, err := os.Open(srcFolder + "/" + fileName + ".xml")
		printError(err)
		defer srcFile.Close()

		data, err := ioutil.ReadAll(srcFile)
		printError(err)

		swsSrcContent := new(SwsStruct)
		err = xml.Unmarshal(data, swsSrcContent)
		if err != nil {
			printError(err)
			continue
		}

		// get column information
		_, err = cfg.GetSection(swsSrcContent.Info.Column)
		if err == nil {
			columnSettings = cfg.Section(swsSrcContent.Info.Column).Keys()
		} else {
			columnSettings = cfg.Section("defaultcolumn").Keys()
			log.Println(swsSrcContent.Info.Column + " section does not exist, use default column settings.")
		}

		// get file modified time
		if swsSrcContent.Info.UpdateTime == "" {
			swsSrcContent.Info.UpdateTime = srcFileInfo.ModTime().Format("2006-01-02")
		}

		swsFileInfo, err := os.Stat(swsFolder + "/" + fileName + ".pdf")
		if os.IsNotExist(err) || srcFileInfo.ModTime().After(swsFileInfo.ModTime()) {
			err = os.Remove(foFolder + "/" + fileName + ".fo")
			printError(err)
			err = os.Remove(swsFolder + "/" + fileName + ".pdf")
			if err != nil && !os.IsNotExist(err) {
				printError(err)
				continue
			}

			cacheFile, err := os.OpenFile(foFolder+"/"+fileName+".fo", os.O_CREATE|os.O_RDWR, os.ModePerm)
			printError(err)
			defer cacheFile.Close()

			foString := foXmlAndRootHead() +
				foLayout() +
				foStaticContent(swsSrcContent) +
				foTableHeadAndColumn() +
				foTableHeaderAndFooter(swsSrcContent.Operator) +
				foTableBody(swsSrcContent) +
				foXmlEnd()

			cacheFile.WriteString(foString)

			pathSeparator := string(os.PathSeparator)
			var fopCommand string = "fop"
			if runtime.GOOS == "windows" {
				fopCommand += ".cmd"
			}
			out, err := exec.Command("fop"+pathSeparator+fopCommand,
				"-c", "fop"+pathSeparator+"fop.xconf",
				"-fo", foFolder+pathSeparator+fileName+".fo",
				"-pdf", swsFolder+pathSeparator+fileName+".pdf").Output()

			printError(err)
			if string(out) != "" {
				log.Println(string(out))
			}
			if err == nil {
				log.Println(fileName + ".pdf is created.")
			} else {
				log.Println(fileName + ".pdf is not created.")
			}

		} else {
			log.Println(srcFileInfo.Name() + " is not changed.")
		}
	}

}

func printError(err error) {
	if err == nil || os.IsNotExist(err) {
		return
	}
	log.Println(err)
}
