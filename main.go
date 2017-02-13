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

	timeFile, err := os.Create("time.txt")
	printError(err)
	timeFile.Close()

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

	var tvgTimeList string
	// main loop, deal with all src xml files
	for _, srcFileInfo := range srcFileInfoList {
		if srcFileInfo.IsDir() {
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

		tvgTimeList += swsSrcContent.Operator.Station +
			swsSrcContent.Operator.Position +
			" " + tvgTime(swsSrcContent) + "\n"

		swsFileInfo, err := os.Stat(swsFolder + "/" + fileName + ".pdf")
		if os.IsNotExist(err) || srcFileInfo.ModTime().After(swsFileInfo.ModTime()) {
			err = os.Remove(foFolder + "/" + fileName + ".fo")
			printError(err)
			err = os.Remove(swsFolder + "/" + fileName + ".pdf")
			if err != nil && !os.IsNotExist(err) {
				printError(err)
				continue
			}

			// get column information
			_, err = cfg.GetSection(swsSrcContent.Info.Column)
			if err == nil {
				columnSettings = cfg.Section(swsSrcContent.Info.Column).Keys()
			} else {
				columnSettings = cfg.Section("defaultcolumn").Keys()
				log.Println(swsSrcContent.Info.Column + " section does not exist, use defaultcolumn settings.")
			}

			// get file modified time
			if swsSrcContent.Info.UpdateTime == "" {
				swsSrcContent.Info.UpdateTime = srcFileInfo.ModTime().Format("2006-01-02")
			}

			cacheFile, err := os.OpenFile(foFolder+"/"+fileName+".fo", os.O_CREATE|os.O_RDWR, os.ModePerm)
			printError(err)
			defer cacheFile.Close()

			foString := foContentString(swsSrcContent)

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

			if err == nil {
				log.Println(fileName + ".pdf is created.")
			} else {
				log.Println(fileName + ".pdf is not created.")
			}
			printError(err)
			if string(out) != "" {
				log.Println(string(out))
			}
		} else {
			log.Println(srcFileInfo.Name() + " is not changed.")
		}
	}

	timeFile, err := os.OpenFile("time.txt", os.O_RDWR|os.O_APPEND, 0666)
	printError(err)
	defer timeFile.Close()
	timeFile.WriteString(tvgTimeList)

	for _, srcFileInfo := range srcFileInfoList {
		if srcFileInfo.IsDir() && !strings.HasPrefix(srcFileInfo.Name(), ".") {
			createSws(srcFolder+"/"+srcFileInfo.Name(), swsFolder+"/"+srcFileInfo.Name())
		}
	}
}

func tvgTime(swsSrcContent *SwsStruct) (totalTime string) {
	var time float32
	for _, process := range swsSrcContent.Operator.Processes {
		time += process.Time
		for _, subprocess := range process.SubProcesses {
			time += subprocess.Time
		}
	}
	totalTime = fmt.Sprintf("%.1f", time*60)

	return
}

func printError(err error) {
	if err == nil || os.IsNotExist(err) {
		return
	}
	log.Println(err)
}
