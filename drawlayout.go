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
	canvas.Scale(0.02)
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
	w := station.Position.W
	h := station.Position.H
	vw := 4600 // vehicle width
	vh := 2000 // vehicle height

	canvas.TranslateRotate(x+w/2, y+h/2, r)

	// station frame
	style := "fill:WhiteSmoke;stroke:black;stroke-width:" + strconv.Itoa(1*h/100)
	canvas.CenterRect(0, 0, w, h, style)

	canvas.Image(-vw/2, -vh/2, vw, vh, "pic/vehicle.svg")
	sfontSize := h / 10
	sfontStyle := "fill:white;font-size:" + strconv.Itoa(sfontSize) + "px"
	var spx, spy, epx, epy int
	ohl, ovl := 14*w/100, 16*h/100 // operator horizontal length and vertical length
	ow := 10 * h / 100             // operator width
	for _, sws := range station.Swses {
		operator := sws.Operator
		Time := totalProcessTime(sws)
		rate := Time.TotalTime / 188
		var strokeColor string
		switch {
		case rate > 0 && rate <= 65:
			strokeColor = "red"
		case rate > 65 && rate <= 85:
			strokeColor = "yellow"
		case rate > 85 && rate <= 95:
			strokeColor = "red"
		case rate > 95:
			strokeColor = "purple"
		default:
			strokeColor = "green"
		}
		lineStyle := "stroke-linecap:round" +
			";stroke:" + strokeColor +
			";stroke-width:" + strconv.Itoa(ow)

		rateString := fmt.Sprintf("%.0f%%", rate)
		switch operator.Position {
		case "RF":
			epx = -5 * w / 100
			spx = epx - ohl
			spy = -1*vh/2 - 4*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, rateString, sfontStyle)
		case "RM":
			epx = ohl / 2
			spx = -1 * epx
			spy = -1*vh/2 - 14*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, rateString, sfontStyle)
		case "RB":
			spx = 5 * w / 100
			epx = spx + ohl
			spy = -1*vh/2 - 4*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, rateString, sfontStyle)
		case "LF":
			epx = -5 * w / 100
			spx = epx - ohl
			spy = vh/2 + 4*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, rateString, sfontStyle)
		case "LM":
			epx = ohl / 2
			spx = -1 * epx
			spy = vh/2 + 14*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, rateString, sfontStyle)
		case "LB":
			spx = 5 * w / 100
			epx = spx + ohl
			spy = vh/2 + 4*h/100
			epy = spy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Text(spx, spy+ow/3, rateString, sfontStyle)

		case "MF":
			spx = -5 * w / 100
			epx = spx
			epy = ovl / 2
			spy = -1 * epy
			canvas.Line(spx, spy, epx, epy, lineStyle)
			canvas.Gtransform("rotate(-90," +
				strconv.Itoa(epx+ow/3) + "," +
				strconv.Itoa(epy) + ")")
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
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
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
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
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
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
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
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
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
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
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
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
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
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
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
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
			canvas.Text(epx+ow/3, epy, rateString, sfontStyle)
			canvas.Gend()
			canvas.Gend()

		}
	}
	// canvas.Gend()
	// canvas.TranslateRotate(x+w/2, y+h/2, r)
	lfontSize := h / 6
	lfontStyle := "text-anchor:middle;fill:white;font-size:" + strconv.Itoa(lfontSize) + "px"
	if r > 135 {
		canvas.Rotate(180)
	}
	canvas.Rect(-w/2, h/2, w, 12*lfontSize/10, style+";fill:Silver")
	canvas.Text(0, h/2+lfontSize, station.Name, lfontStyle)
	if r > 135 {
		canvas.Gend()
	}

	canvas.Gend()

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
