package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/svgo"
)

func drawLayout(area AreaStruct) {
	svgFileName := area.Name + ".svg"
	svgFile, err := os.Create("layout/" + svgFileName)
	printError(err)
	defer svgFile.Close()

	width := 1000
	height := 500
	canvas := svg.New(svgFile)
	canvas.Start(width, height)
	canvas.Scale(0.01)
	canvas.Rect(0, 0, width, height, "fill:none;stroke:black")
	// canvas.Circle(width/2, height/2, 100)
	// canvas.Image(width/4, height/4, width/2, height/2, "pic/a.png")
	for _, station := range area.Stations {
		drawStation(canvas, station)
	}
	canvas.Gend()
	canvas.End()
}

func drawStation(canvas *svg.SVG, station StationStruct) {
	x := station.Position.X
	y := station.Position.Y
	r := station.Position.R
	vw := 4600 // vehicle width
	vh := 2000 // vehicle height
	w := 6500
	h := 4500

	canvas.TranslateRotate(x+w/2, y+h/2, r)

	// station frame
	style := "fill:none;stroke:black;stroke-linecap:round;stroke-width:" + strconv.Itoa(1*h/100)
	canvas.CenterRect(0, 0, w, h, style)

	canvas.Image(-vw/2, -vh/2, vw, vh, "pic/vehicle.svg")
	sfontSize := h / 10
	sfontStyle := "fill:blue;font-size:" + strconv.Itoa(sfontSize) + "px"
	var spx, spy, epx, epy int
	ohl, ovl := 14*w/100, 16*h/100 // operator horizontal length and vertical length
	ow := 10 * h / 100             // operator width
	for _, sws := range station.Swses {
		operator := sws.Operator
		lineStyle := "fill:none;stroke:black;stroke-linecap:round;stroke-width:" + strconv.Itoa(ow)
		switch operator.Position {
		case "RF":
			epx = -5 * w / 100
			spx = epx - ohl
			spy = -1*vh/2 - 4*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, station.Name, sfontStyle)
		case "RM":
			epx = ohl / 2
			spx = -1 * epx
			spy = -1*vh/2 - 14*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, station.Name, sfontStyle)
		case "RB":
			spx = 5 * w / 100
			epx = spx + ohl
			spy = -1*vh/2 - 4*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, station.Name, sfontStyle)
		case "LF":
			epx = -5 * w / 100
			spx = epx - ohl
			spy = vh/2 + 4*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, station.Name, sfontStyle)
		case "LM":
			epx = ohl / 2
			spx = -1 * epx
			spy = vh/2 + 14*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, station.Name, sfontStyle)
		case "LB":
			spx = 5 * w / 100
			epx = spx + ohl
			spy = vh/2 + 4*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, station.Name, sfontStyle)

		case "MF":
			spx = -5 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()
		case "MM":
			spx = 5 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()
		case "MB":
			spx = 15 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()

		case "FM":
			spx = -38 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()
		case "FR":
			spx = -30 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Rotate(-20)
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()
			canvas.Gend()
		case "FL":
			spx = -30 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Rotate(20)
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()
			canvas.Gend()

		case "BM":
			spx = 38 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()
		case "BR":
			spx = 30 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Rotate(-20)
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()
			canvas.Gend()
		case "BL":
			spx = 30 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Rotate(20)
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, station.Name, sfontStyle)
			canvas.Gend()
			canvas.Gend()

		}
	}
	// canvas.Rect(0, 0, w, h, "fill:black")
	lfontSize := h / 6
	lfontStyle := "text-anchor:middle;fill:blue;font-size:" + strconv.Itoa(lfontSize) + "px"
	canvas.Text(0, lfontSize/3, station.Name, lfontStyle)

	canvas.Gend()

	// canvas.Gend()
	return
}

func createLayout(srcFolder, swsFolder string) {
	os.MkdirAll(swsFolder, os.ModeDir|os.ModePerm)

	fmt.Printf("\n-dir %s\n", srcFolder)
	srcFileInfoList, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		printError(err)
		return
	}

	// timeList := "\"" + srcFolder + "\"\n"
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
	}
}
