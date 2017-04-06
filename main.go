package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-ini/ini"
)

var cfg *ini.File
var columnSettings []*ini.Key

var foFolder string = "fo"
var timeFileName string = "time.csv"
var stations map[string]StationStruct

type StationStruct struct {
	Name          string
	OperatorInfos []OperatorInfoStruct
}

type OperatorInfoStruct struct {
	StationName   string
	Position      string
	OperationTime OperationTimeStruct
	SwsContent    SwsStruct
}

func main() {
	stations = make(map[string]StationStruct)

	var err error
	cfg, err = ini.Load("sws.ini")
	printError(err)
	srcFolder := cfg.Section("general").Key("srcfolder").MustString("src")
	swsFolder := cfg.Section("general").Key("swsfolder").MustString("sws")

	timeFile, err := os.Create(timeFileName)
	printError(err)
	csvFileTitle := "station,valueTime,noneValueTime,waitingTime,totalTime\n"
	timeFile.WriteString(csvFileTitle)
	timeFile.Close()

	os.MkdirAll(foFolder, os.ModeDir|os.ModePerm)
	createSws(srcFolder, swsFolder)
	os.RemoveAll(foFolder)

	// GenerateXslFile(stations)
	// fmt.Printf("%v\n", stations)
}

func createSws(srcFolder, swsFolder string) {
	os.MkdirAll(swsFolder, os.ModeDir|os.ModePerm)

	fmt.Printf("\n-dir %s\n", srcFolder)
	srcFileInfoList, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		printError(err)
		return
	}

	tvgTimeList := srcFolder + "\n"
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

		// swsSrcContent := new(SwsStruct)
		var swsSrcContent SwsStruct
		err = xml.Unmarshal(data, &swsSrcContent)
		if err != nil {
			printError(err)
			continue
		}

		// output time information to time.csv
		operationTime := totalProcessTime(swsSrcContent)
		tvgTimeList += swsSrcContent.Operator.Station +
			"_" + swsSrcContent.Operator.Position + "," +
			fmt.Sprintf("%.1f,", operationTime.ValueTime) +
			fmt.Sprintf("%.1f,", operationTime.NoneValueTime) +
			fmt.Sprintf("%.1f,", operationTime.WaitingTime) +
			fmt.Sprintf("%.1f,", operationTime.TotalTime) +
			"\n"

		// compare src and pdf modified time. if src file new, create pdf
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
			if swsSrcContent.Info.Column == "" {
				swsSrcContent.Info.Column = "defaultcolumn"
			}
			columnSetion, err := cfg.GetSection(swsSrcContent.Info.Column)
			if err != nil {
				columnSetion = cfg.Section("defaultcolumn")
				log.Println(swsSrcContent.Info.Column + " section does not exist, use defaultcolumn settings.")
			}
			columnSettings = columnSetion.Keys()

			// get file modified time
			if swsSrcContent.Info.UpdateTime == "" {
				swsSrcContent.Info.UpdateTime = srcFileInfo.ModTime().Format("2006-01-02")
			}

			// build fo temperory file
			cacheFile, err := os.OpenFile(foFolder+"/"+fileName+".fo", os.O_CREATE|os.O_RDWR, os.ModePerm)
			printError(err)
			defer cacheFile.Close()
			foString := foContentString(swsSrcContent)
			cacheFile.WriteString(foString)

			pathSeparator := string(os.PathSeparator)
			var fopCommand string = "fop"
			// if runtime.GOOS == "windows" {
			// fopCommand += ".cmd"
			// }
			out, err := exec.Command("fop"+pathSeparator+fopCommand,
				"-c", "fop"+pathSeparator+"fop.xconf",
				"-fo", foFolder+pathSeparator+fileName+".fo",
				"-pdf", swsFolder+pathSeparator+fileName+".pdf").Output()

			if err == nil {
				log.Println(fileName + ".pdf is created.")
			} else {
				log.Println(fileName + ".pdf is not created, error information is :")
			}
			printError(err)
			if string(out) != "" {
				log.Println(string(out))
			}
		} else {
			log.Println(srcFileInfo.Name() + " is not changed.")
		}

		fillOperatorInfoToStation(swsSrcContent)
	}

	// write time information to time.csv file
	timeFile, err := os.OpenFile(timeFileName, os.O_RDWR|os.O_APPEND, 0666)
	printError(err)
	defer timeFile.Close()
	tvgTimeList += "\n"
	timeFile.WriteString(tvgTimeList)

	// build sub dir files
	for _, srcFileInfo := range srcFileInfoList {
		if srcFileInfo.IsDir() && !strings.HasPrefix(srcFileInfo.Name(), ".") {
			createSws(srcFolder+"/"+srcFileInfo.Name(), swsFolder+"/"+srcFileInfo.Name())
		}
	}
}

func fillOperatorInfoToStation(swsSrcContent SwsStruct) {
	stationName := swsSrcContent.Operator.Station
	position := swsSrcContent.Operator.Position
	var station StationStruct
	station.Name = stationName
	_, isPresent := stations[stationName]
	if !isPresent {
		stations[stationName] = station
	} else {
		var isPositionPresent bool = false
		s := stations[stationName]
		for _, o := range s.OperatorInfos {
			if o.Position == position {
				isPositionPresent = true
				break
			}
		}
		if isPositionPresent == true {
			return
		}
	}

	var operatorInfo OperatorInfoStruct
	operatorInfo.StationName = stationName
	operatorInfo.Position = swsSrcContent.Operator.Position
	operatorInfo.OperationTime = totalProcessTime(swsSrcContent)
	operatorInfo.SwsContent = swsSrcContent

	station.OperatorInfos = append(stations[stationName].OperatorInfos, operatorInfo)
	stations[stationName] = station
	// fmt.Println(station.Name + " " + stationName)
	// fmt.Printf("%#v\n", areas[i].Stations[j])
}

func printError(err error) {
	if err == nil || os.IsNotExist(err) {
		return
	}
	log.Println(err)
}
