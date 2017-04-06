package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var xslFileName string = "L.xsl"

func GenerateXslFile(stations map[string]StationStruct) {
	readStationsInfoFromCsvFile()
	xslFile, err := os.Create(xslFileName)
	printError(err)
	defer xslFile.Close()

	tact := cfg.Section("general").Key("tact").MustFloat64(188)
	xsl := NewXSL(xslFile)
	for _, station := range stations {
		for _, operator := range station.OperatorInfos {
			xsl.AddOperator(operator, tact)
		}
	}
	xsl.End()
}

type XSL struct {
	Writer io.Writer
}

func NewXSL(w io.Writer) *XSL {
	x := new(XSL)
	x.Writer = w

	x.Writer.Write(xslHead)
	return x
}

func (self *XSL) AddOperator(operatorInfo OperatorInfoStruct, tact float64) (err error) {
	stationName := strings.TrimSpace(operatorInfo.StationName)
	begin := 0
	end := len(stationName)
	if end > 5 {
		begin = end - 5
	}
	stationName = stationName[begin:end]

	position := operatorInfo.Position
	totalTime := operatorInfo.OperationTime.TotalTime
	workload := totalTime / tact * 100
	var positionColourCode string
	switch {
	case workload >= 0 && workload <= 70:
		// positionColour = "red"
		positionColourCode = "D40000"
	case workload > 70 && workload <= 80:
		// positionColour = "yellow"
		positionColourCode = "D4AA00"
	case workload > 80 && workload <= 95:
		// positionColour = "green"
		positionColourCode = "008000"
	case workload > 95:
		// positionColour = "purple"
		positionColourCode = "800080"
	default:
		// positionColour = "green"
		positionColourCode = "008000"
	}

	template := `
	<xsl:template match="node()[@id='` + stationName + `']/svg:g[@inkscape:label='operator']/svg:image[@inkscape:label='` + position + `']">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<xsl:attribute name="xlink:href">
				<xsl:text>pic/colour/` + positionColourCode + `.svg</xsl:text>
			</xsl:attribute>
		</xsl:copy>
	</xsl:template>
	<xsl:template match="node()[@id='` + stationName + `']/svg:g[@inkscape:label='workload']/svg:text[@inkscape:label='` + position + `']">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<xsl:text>` + fmt.Sprintf(position+":%.0f%%", workload) + `</xsl:text>
		</xsl:copy>
	</xsl:template>

	<xsl:template match="node()[@id='c` + stationName + `']/svg:g[@inkscape:label='operator']/svg:rect[@inkscape:label='` + position + `']">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<xsl:attribute name="style">
				<xsl:call-template name="replace-string">
					<xsl:with-param name="text" select="@style"/>
					<xsl:with-param name="replace" select="'#d35f5f'" />
					<xsl:with-param name="with" select="'#` + positionColourCode + `'"/>
				</xsl:call-template>
			</xsl:attribute>
		</xsl:copy>
	</xsl:template>
	<xsl:template match="node()[@id='c` + stationName + `']/svg:g[@inkscape:label='workload']/svg:text[@inkscape:label='` + position + `']">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<xsl:text>` + fmt.Sprintf("%.0f / %.0f%%", totalTime, workload) + `</xsl:text>
		</xsl:copy>
	</xsl:template>
`

	_, err = self.Writer.Write([]byte(template))
	return
}

func (self *XSL) AddTemplate(path, attributeName, attributeValue, text string) (err error) {
	template := []byte(`	<xsl:template match="` + path + `">
		<xsl:copy>
			<xsl:apply-templates select="@*[not(local-name()='` + attributeName + `')]"/>
			<xsl:attribute name="` + attributeName + `">
				<xsl:text>` + attributeValue + `</xsl:text>
			</xsl:attribute>
			<xsl:apply-templates select="node()[not(.)]"/>
			<xsl:text>` + text + `</xsl:text>
		</xsl:copy>
	</xsl:template>
`)

	_, err = self.Writer.Write(template)
	return
}

func (self *XSL) End() (err error) {
	_, err = self.Writer.Write(xslEnd)
	return
}

func readStationsInfoFromCsvFile() {
	file, err := os.Open("stations.csv")
	printError(err)
	defer file.Close()

	// var buffer bytes.Buffer
	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		info := strings.Split(string(lineBytes), ",")

		stationName := info[0]
		var station StationStruct
		station.Name = stationName
		_, isPresent := stations[stationName]
		if !isPresent {
			stations[stationName] = station
		} else {
			var isPositionPresent bool = false
			s := stations[stationName]
			for _, o := range s.OperatorInfos {
				if o.Position == info[1] {
					isPositionPresent = true
					break
				}
			}
			if isPositionPresent == true {
				continue
			}
		}

		var operatorInfo OperatorInfoStruct
		operatorInfo.StationName = info[0]
		operatorInfo.Position = info[1]
		operatorInfo.OperationTime.TotalTime, _ = strconv.ParseFloat(info[2], 64)

		station.OperatorInfos = append(stations[stationName].OperatorInfos, operatorInfo)
		stations[stationName] = station
	}
}

var xslHead []byte = []byte(`<xsl:stylesheet version="1.0" 
	xmlns:svg="http://www.w3.org/2000/svg"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
	xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
	xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
	<xsl:output method="xml" version="1.0" encoding="UTF-8" indent="yes"/>
	<xsl:strip-space elements="*"/>

	<!-- <xsl:key name="id" match="ID" use="." /> -->
	<xsl:template match="@*|node()">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
		</xsl:copy>
	</xsl:template>

	<xsl:template name="replace-string">
		<xsl:param name="text"/>
		<xsl:param name="replace"/>
		<xsl:param name="with"/>
		<xsl:choose>
			<xsl:when test="contains($text,$replace)">
				<xsl:value-of select="substring-before($text,$replace)"/>
				<xsl:value-of select="$with"/>
				<xsl:call-template name="replace-string">
					<xsl:with-param name="text" select="substring-after($text,$replace)"/>
					<xsl:with-param name="replace" select="$replace"/>
					<xsl:with-param name="with" select="$with"/>
				</xsl:call-template>
			</xsl:when>
			<xsl:otherwise>
				<xsl:value-of select="$text"/>
			</xsl:otherwise>
		</xsl:choose>
	</xsl:template>

	<xsl:template match="svg:g[@inkscape:label='operator']/svg:image">
		<xsl:apply-templates />
	</xsl:template>
	<xsl:template match="svg:g[@inkscape:label='operator']/svg:rect">
		<xsl:apply-templates />
	</xsl:template>
	<xsl:template match="svg:g[@inkscape:label='workload']//node()">
		<xsl:apply-templates />
	</xsl:template>
`)
var xslEnd []byte = []byte(`</xsl:stylesheet>`)
