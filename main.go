package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

var cfg *ini.File

// var layoutCfg *ini.File
var columnSettings []*ini.Key
var foFolder string = "fo"
var timeFileName string = "time.csv"
var areas []AreaStruct

func main() {
	var err error
	cfg, err = ini.Load("sws.ini")
	printError(err)
	srcFolder := cfg.Section("general").Key("srcfolder").MustString("src")
	swsFolder := cfg.Section("general").Key("swsfolder").MustString("sws")

	initAreas()

	timeFile, err := os.Create(timeFileName)
	printError(err)
	csvFileTitle := "station,valueTime,noneValueTime,waitingTime,totalTime\n"
	timeFile.WriteString(csvFileTitle)
	timeFile.Close()

	os.MkdirAll(foFolder, os.ModeDir|os.ModePerm)
	createSws(srcFolder, swsFolder)
	os.RemoveAll(foFolder)

	// fmt.Printf("%v\n", areas)
	for _, area := range areas {
		drawLayout(area)
	}
}

func initAreas() {
	layoutCfg, err := ini.Load("layout.ini")
	printError(err)
	sections := layoutCfg.Sections()
	for _, section := range sections {
		if section.Name() == "DEFAULT" {
			continue
		}

		var area AreaStruct
		area.Name = section.Name()
		keys := section.Keys()
		for _, key := range keys {
			position := key.Strings(",")
			switch key.Name() {
			case "position":
				area.Position.X, _ = strconv.Atoi(position[0])
				area.Position.Y, _ = strconv.Atoi(position[1])
				area.Position.R, _ = strconv.ParseFloat(position[2], 64)
				area.Position.W, _ = strconv.Atoi(position[3])
				area.Position.H, _ = strconv.Atoi(position[4])
				continue

			case "papersize":
				area.Paper.W, _ = strconv.Atoi(position[0])
				area.Paper.H, _ = strconv.Atoi(position[1])
				continue

			default:
				var station StationStruct
				station.Name = key.Name()
				station.Position.X, _ = strconv.Atoi(position[0])
				station.Position.Y, _ = strconv.Atoi(position[1])
				station.Position.R, _ = strconv.ParseFloat(position[2], 64)
				station.Position.W, _ = strconv.Atoi(position[3])
				station.Position.H, _ = strconv.Atoi(position[4])
				station.Position.VW, _ = strconv.Atoi(position[5])
				station.Position.VH, _ = strconv.Atoi(position[6])
				station.Position.Kind = position[7]
				area.Stations = append(area.Stations, station)
			}
		}
		areas = append(areas, area)
	}
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
	for i, area := range areas {
		for j, station := range area.Stations {
			if station.Name != stationName {
				continue
			}
			// var operatorInfo OperatorInfoStruct
			// operatorInfo.Position = swsSrcContent.Operator.Position
			// operatorInfo.OperationTime = totalProcessTime(swsSrcContent)
			areas[i].Stations[j].Swses = append(areas[i].Stations[j].Swses, swsSrcContent)
			// fmt.Println(station.Name + " " + stationName)
			// fmt.Printf("%#v\n", areas[i].Stations[j])
		}
	}
}

func printError(err error) {
	if err == nil || os.IsNotExist(err) {
		return
	}
	log.Println(err)
}
