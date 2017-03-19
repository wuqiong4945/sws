package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type XSL struct {
	Writer io.Writer
}

func NewXSL(w io.Writer) *XSL {
	x := new(XSL)
	x.Writer = w

	x.Writer.Write(xslHead)
	return x
}

func (self *XSL) AddOperator(station, position string, operatorInfo OperatorInfoStruct) (err error) {
	template := `
	<xsl:template match="node()[@id='60001']/svg:g[@inkscape:label='operator']/svg:image[@inkscape:label='RM']">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<xsl:attribute name="xlink:href">
				<xsl:text>pic/purple.png</xsl:text>
			</xsl:attribute>
		</xsl:copy>
	</xsl:template>

	<xsl:template match="node()[@id='60001']/svg:g[@inkscape:label='workload']/svg:text[@inkscape:label='RM']">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<xsl:text>added</xsl:text>
		</xsl:copy>
	</xsl:template>


	<xsl:template match="node()[@id='c60001']/svg:g[@inkscape:label='operator']/svg:rect[@inkscape:label='RM']">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<xsl:attribute name="style">
				<xsl:call-template name="replace-string">
					<xsl:with-param name="text" select="@style"/>
					<xsl:with-param name="replace" select="'#d35f5f'" />
					<xsl:with-param name="with" select="'blue'"/>
				</xsl:call-template>
			</xsl:attribute>
		</xsl:copy>
	</xsl:template>

	<xsl:template match="node()[@id='c60001']/svg:g[@inkscape:label='workload']/svg:text[@inkscape:label='RM']">
		<xsl:copy>
			<xsl:apply-templates select="@*|node()"/>
			<xsl:text>added</xsl:text>
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
