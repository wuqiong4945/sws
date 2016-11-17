package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"
)

func foXmlAndRootHead() string {
	fontfamily := cfg.Section("font").Key("fontfamily").MustString("Microsoft JhengHei, serif")
	fontsize := cfg.Section("font").Key("fontsize").MustString("7")

	foXmlHead := `<?xml version="1.0" encoding="utf-8"?>` + "\n"
	foRootHead := `<fo:root xmlns:fo="http://www.w3.org/1999/XSL/Format"` +
		` font-family="` + fontfamily + `" font-size="` + fontsize + `"` +
		` font-selection-strategy="character-by-character" language="en">
		`
	return foXmlHead + foRootHead
}

func foLayout() string {
	s := `
	<fo:layout-master-set>
    <fo:simple-page-master master-name="A4" page-width="297mm" page-height="210mm" margin-top="10mm" margin-bottom="10mm" margin-left="1mm" margin-right="1mm">
      <fo:region-body margin-top="10mm" margin-bottom="20mm" margin-left="1mm" margin-right="1mm" />
      <fo:region-before region-name="page-head" extent="10mm"/>
      <fo:region-after region-name="page-foot" extent="20mm"/>
      <fo:region-start region-name="page-start" extent="1mm"/>
      <fo:region-end region-name="page-end" extent="1mm"/>
    </fo:simple-page-master>

    <fo:page-sequence-master master-name="standard">
      <fo:repeatable-page-master-reference  master-reference="A4"/>        
      <!-- <fo:repeatable-page-master-alternatives> -->
        <!-- <fo:conditional-page-master-reference master-reference="first" page-position="first" /> -->
        <!-- <fo:conditional-page-master-reference master-reference="left" odd-or-even="even" /> -->
        <!-- <fo:conditional-page-master-reference master-reference="right" odd-or-even="odd" /> -->
        <!-- </fo:repeatable-page-master-alternatives> -->
    </fo:page-sequence-master>
  </fo:layout-master-set>
	`
	return s
}

func foStaticContent(swsSrcContent *SwsStruct) string {
	info := swsSrcContent.Info
	safety := swsSrcContent.Operator.Safety

	var title string
	if info.Title == "" {
		title = "标 准 工 艺 操 作 指 导"
	} else {
		title = info.Title
	}

	s := `
	<fo:page-sequence master-reference="standard">
    <!-- page head -->
    <fo:static-content flow-name="page-head">
      <fo:block>
			<fo:table display-align="center" border-collapse="collapse" table-layout="fixed" width="100%">
				<fo:table-column column-width="30mm"/>
				<fo:table-column column-width="233mm"/>
				<fo:table-column column-width="30mm"/>

        <fo:table-body>
          <fo:table-row height="10mm" border-color="black" border-width="0.75pt" border-style="solid">
            <fo:table-cell><fo:block/></fo:table-cell>
            <fo:table-cell><fo:block text-align="center" font-size="15pt">` + title + `</fo:block></fo:table-cell>
            <fo:table-cell>
							<fo:block text-align="center">
								<fo:external-graphic content-width="9mm" content-height="9mm" scaling="non-uniform" src="fop/images/bmw.png" />
							</fo:block>
						</fo:table-cell>
          </fo:table-row>
				</fo:table-body>

        </fo:table>
      </fo:block>
    </fo:static-content>

    <!-- page foot -->
    <fo:static-content flow-name="page-foot">
      <fo:block text-align="right" vertical-align="middle">

        <fo:table font-size="8" display-align="center" border-collapse="collapse" table-layout="fixed" width="100%">
          <fo:table-column       column-width= " 30mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
          <fo:table-column       column-width= " 60mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
          <fo:table-column       column-width= " 1mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
          <fo:table-column       column-width= " 85mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
          <fo:table-column       column-width= " 117mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />

          <fo:table-body>
            <fo:table-row height="5mm" border-color="black" border-width="0.75pt" border-style="solid">
              <fo:table-cell><fo:block text-align="center">创建</fo:block></fo:table-cell>
              <fo:table-cell><fo:block text-align="center">` + info.Author + `</fo:block></fo:table-cell>

              <fo:table-cell empty-cells="show" number-rows-spanned="4"> <fo:block/> </fo:table-cell>

              <fo:table-cell number-rows-spanned="4" border-end-color="white" border-width="0.75pt" border-style="solid">
                <fo:block text-align="left">
                  <fo:instream-foreign-object>
                    <svg:svg xmlns:svg="http://www.w3.org/2000/svg" width="20px" height="10px">
                      <svg:g style="fill:red; stroke:#000000">
                        <svg:rect x="10" y="0" width="10" height="10" />
                      </svg:g>
                    </svg:svg>
                  </fo:instream-foreign-object>
                  关键工序
                </fo:block>
              </fo:table-cell>

              <fo:table-cell number-rows-spanned="4" border-start-color="white" border-width="0.75pt" border-style="solid">
                <fo:block text-align="right">
								`
	if safety.IsESDShoes == "yes" {
		s += `<fo:external-graphic content-width="23mm" content-height="20mm" scaling="non-uniform" src="fop/images/shoe.png"/>` + "\n"
	}
	if safety.IsWorkware == "yes" {
		s += `<fo:external-graphic content-width="23mm" content-height="20mm" scaling="non-uniform" src="fop/images/clothes.png"/>` + "\n"
	}
	if safety.IsSafetyGlasses == "yes" {
		s += `<fo:external-graphic content-width="23mm" content-height="20mm" scaling="non-uniform" src="fop/images/glasses.png"/>` + "\n"
	}
	if safety.IsSafetyGloves == "yes" {
		s += `<fo:external-graphic content-width="23mm" content-height="20mm" scaling="non-uniform" src="fop/images/glove.png"/>` + "\n"
	}

	s += `        </fo:block>
              </fo:table-cell>

            </fo:table-row>
            <fo:table-row height="5mm" border-color="black" border-width="0.75pt" border-style="solid">
              <fo:table-cell><fo:block text-align="center">批准</fo:block></fo:table-cell>
              <fo:table-cell><fo:block/></fo:table-cell>
            </fo:table-row>
            <fo:table-row height="5mm" border-color="black" border-width="0.75pt" border-style="solid">
              <fo:table-cell><fo:block text-align="center">部门</fo:block></fo:table-cell>
              <fo:table-cell><fo:block text-align="center">` + info.Department + `</fo:block></fo:table-cell>
            </fo:table-row>
            <fo:table-row height="5mm" border-color="black" border-width="0.75pt" border-style="solid">
              <fo:table-cell><fo:block text-align="center">更新</fo:block></fo:table-cell>
              <fo:table-cell><fo:block text-align="center">` + info.UpdateTime + `</fo:block></fo:table-cell>
            </fo:table-row>
          </fo:table-body>
        </fo:table>

      </fo:block>

      <!--<fo:block font-size="10pt" text-align="end">Page <fo:page-number/> of <fo:page-number-citation ref-id="TheVeryLastPage"/></fo:block>-->
    </fo:static-content>
	`
	return s
}

func foTableHeadAndColumn() string {
	foTableHead := `
    <fo:flow flow-name="xsl-region-body">
      <fo:table display-align="center" border-collapse="collapse" table-layout="fixed" width="100%">
			`
	foTableColumnPic := `
        <fo:table-column id= " model     " column-width= " 15mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
        <fo:table-column id= " model1    " column-width= " 15mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
        <fo:table-column id= " station   " column-width= " 15mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
        <fo:table-column id= " station1  " column-width= " 15mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
        <fo:table-column id= " operator  " column-width= " 15mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
        <fo:table-column id= " operator1 " column-width= " 15mm " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
        <fo:table-column id= " break     " column-width= " 1mm  " border-color= " black " border-width= " 0.75pt " border-style= " solid " />
				`
	var foTableColumnText string
	for _, columnSetting := range columnSettings {
		vals := columnSetting.Strings(",")
		columnTextString := `<fo:table-column id= "` + vals[0] + `" column-width= "` + vals[1] + `mm" border-color= "black" border-width= "0.75pt" border-style= "solid" />`
		foTableColumnText += columnTextString + "\n"
	}

	return foTableHead + foTableColumnPic + foTableColumnText
}

func foTableHeaderAndFooter(operator OperatorStruct) string {
	foTableHeaderPic := `
        <fo:table-header>
          <fo:table-row height="4mm" font-size="7pt" font-weight="bold" background-color="Gainsboro" border-color="black" border-width="0.75pt" border-style="solid">
            <fo:table-cell><fo:block text-align="center">车型</fo:block></fo:table-cell>
            <fo:table-cell><fo:block text-align="center">` + operator.Model + `</fo:block></fo:table-cell>
            <fo:table-cell><fo:block text-align="center">工位</fo:block></fo:table-cell>
            <fo:table-cell><fo:block text-align="center">` + operator.Station + `</fo:block></fo:table-cell>
            <fo:table-cell><fo:block text-align="center">操作者</fo:block></fo:table-cell>
            <fo:table-cell><fo:block text-align="center">` + operator.Position + `</fo:block></fo:table-cell>
						<fo:table-cell background-color="white" border-after-color="white" border-width="0.75pt" border-style="solid"><fo:block/></fo:table-cell>
				`

	var foTableHeaderText string
	for _, columnSetting := range columnSettings {
		vals := columnSetting.Strings(",")
		headString := `<fo:table-cell><fo:block text-align="center">` + vals[2] + `</fo:block></fo:table-cell>`
		foTableHeaderText += headString + "\n"
	}
	foTableHeaderText += `</fo:table-row></fo:table-header>` + "\n"

	foTableFooter := `
        <fo:table-footer>
          <fo:table-row><fo:table-cell><fo:block/></fo:table-cell></fo:table-row>
        </fo:table-footer>
				`

	return foTableHeaderPic + foTableHeaderText + foTableFooter
}

func foTableBody(swsSrcContent *SwsStruct) string {
	xmlPictureCellHead := `
			<fo:table-body>
          <fo:table-row border-color="black" border-width="0.75pt" border-style="solid">
            <fo:table-cell display-align="before" number-columns-spanned="6" number-rows-spanned="152">
						`
	xmlPictureCellEnd := `</fo:table-cell></fo:table-row>` + "\n"

	var xmlTextCellString string
	var processNumber int = swsSrcContent.Operator.FirstProcessNumber
	if processNumber == 0 {
		// if first page, FirstProcessNumber is not set. ProcessNumber stars from 1.
		processNumber = 1
	}
	var processContent []ProcessContent
	for _, process := range swsSrcContent.Operator.Processes {
		content := processTableBodyContent(process, strconv.Itoa(processNumber))
		processContent = append(processContent, content...)
		processNumber++
	}

	var xmlPictureCellBody string
	imagefolder := cfg.Section("general").Key("imagefolder").MustString("images")
	picPositionRightX, picPositionTopY, picPositionDownY := 0, 0, 0
	picWidth, picHeight := 0, 0
	picCellWidth, picCellHeight := 91, 158
	picAspectRatio := "none"
	xmlPictureCellBody += `<fo:block>
		<fo:instream-foreign-object>
			<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"` +
		` width="90mm" height="148mm">
		`
	for _, c := range processContent {
		if c.ProcessPictureName != "" {
			switch c.ProcessPictureSize {
			case "small":
				picWidth, picHeight = 44, 29
			case "long":
				picWidth, picHeight = 89, 37
			case "medium":
				picWidth, picHeight = 89, 44
			case "big":
				picWidth, picHeight = 89, 59
			case "default":
				picWidth, picHeight = 44, 37
			default:
				picWidth, picHeight = 44, 37
			}

			if picPositionRightX+picWidth > picCellWidth {
				picPositionRightX = 0
				picPositionTopY = picPositionDownY
			}
			if picPositionTopY+picHeight > picCellHeight {
				xmlPictureCellBody += `</svg:svg></fo:instream-foreign-object></fo:block>` + "\n"
				xmlPictureCellBody += `<fo:block>
					<fo:instream-foreign-object>
						<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"` +
					` width="90mm" height="148mm">
					`
				picPositionRightX, picPositionTopY, picPositionDownY = 0, 0, 0
			}

			xmlPictureCellBody += `<svg:g>` + "\n"

			imageUrl := imagefolder + "/" + c.ProcessPictureName
			_, err := os.Stat(imageUrl)
			if err == nil {
				// use real image ratio if imagesize is "long"
				if c.ProcessPictureSize == "long" {
					imageFile, _ := os.Open(imageUrl)
					var imageConf image.Config
					switch {
					case strings.HasSuffix(imageUrl, ".png"):
						imageConf, _ = png.DecodeConfig(imageFile)
					case strings.HasSuffix(imageUrl, ".jpg"):
						imageConf, _ = jpeg.DecodeConfig(imageFile)
					default:
						imageConf.Width, imageConf.Height = picWidth, picHeight
					}
					picHeight = imageConf.Height * picWidth / imageConf.Width
				}
				xmlPictureCellBody += `<svg:image` +
					` x="` + strconv.Itoa(picPositionRightX) + `mm"` +
					` y="` + strconv.Itoa(picPositionTopY) + `mm"` +
					` width="` + strconv.Itoa(picWidth) + `mm"` +
					` height="` + strconv.Itoa(picHeight) + `mm"` +
					` preserveAspectRatio="` + picAspectRatio + `"` +
					` xlink:href="` + imageUrl + `" />
					`
			} else {
				// write the url so that we know which pic is missing
				xmlPictureCellBody += `<svg:text` +
					` x="` + strconv.Itoa(picPositionRightX+2) + `mm"` +
					` y="` + strconv.Itoa(picPositionTopY+10) + `mm">` +
					imageUrl +
					`</svg:text>
					`
			}
			// draw the left front number rectangle
			xmlPictureCellBody += `
									<svg:rect` +
				` x="` + strconv.Itoa(picPositionRightX) + `mm"` +
				` y="` + strconv.Itoa(picPositionTopY) + `mm"` +
				` width="25" height="12" style="fill:yellow"/>
									<svg:text` +
				` x="` + strconv.Itoa(picPositionRightX+2) + `mm"` +
				` y="` + strconv.Itoa(picPositionTopY+3) + `mm">` +
				c.ProcessNumber +
				`</svg:text>
				`
			xmlPictureCellBody += `</svg:g>` + "\n"

			picPositionRightX += picWidth + 1
			if picPositionTopY+picHeight+1 > picPositionDownY {
				picPositionDownY = picPositionTopY + picHeight + 1
			}
		}

		xmlTextCellString += c.ProcessTextContent
	}

	xmlPictureCellBody += `</svg:svg></fo:instream-foreign-object></fo:block>` + "\n"
	xmlPictureCellString := xmlPictureCellHead + xmlPictureCellBody + xmlPictureCellEnd

	for i := 0; i < 160; i++ {
		xmlTextCellString += `<fo:table-row height="5mm" border-color="black" border-width="0.75pt" border-style="solid">
        <fo:table-cell border-before-color="white" border-after-color="white" border-width="0.75pt" border-style="solid">
					<fo:block/>
				</fo:table-cell>
				<fo:table-cell><fo:block/></fo:table-cell>
				<fo:table-cell><fo:block/></fo:table-cell>
				<fo:table-cell><fo:block/></fo:table-cell>
				<fo:table-cell><fo:block/></fo:table-cell>
      </fo:table-row>
			`
	}

	xmlTextCellString += `</fo:table-body>` + "\n"
	return xmlPictureCellString + xmlTextCellString
}

type ProcessContent struct {
	ProcessNumber      string
	ProcessPictureName string
	ProcessPictureSize string
	ProcessTextContent string
}

func processTableBodyContent(process ProcessStruct, processNumberString string) (processContent []ProcessContent) {
	var content ProcessContent
	content.ProcessNumber = processNumberString
	content.ProcessPictureName = process.Image
	content.ProcessPictureSize = process.ImageSize

	var processTextContent, backgroundColour, fontWeight string
	// set main item background
	if !strings.Contains(processNumberString, ".") {
		backgroundColour = ""
		fontWeight = ` font-weight="bold"`
	} else {
		// backgroundColour = ` background-color="GhostWhite"`
		backgroundColour = ""
		fontWeight = ``
	}
	processTextContent += `<fo:table-row` + backgroundColour + fontWeight + ` height="5mm" vertical-align="middle" border-color="black" border-width="0.75pt" border-style="solid">` + "\n"

	processTextContent += `<fo:table-cell border-before-color="white" border-after-color="white" border-width="0.75pt" border-style="solid"><fo:block/></fo:table-cell>` + "\n"
	for _, columnSetting := range columnSettings {
		vals := columnSetting.Strings(",")
		switch vals[0] {
		case "number":
			if process.IsKey == "yes" {
				backgroundColour = ` background-color="red"`
			} else {
				backgroundColour = ""
			}
			processTextContent += `<fo:table-cell` + backgroundColour + `><fo:block text-align="left">` + processNumberString + `</fo:block></fo:table-cell>` + "\n"
		case "option":
			op := process.Option
			if op == "" {
				op = "所有配置"
			}
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + op + `</fo:block></fo:table-cell>` + "\n"
		case "tvg":
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + process.Tvg + `</fo:block></fo:table-cell>` + "\n"
		case "description":
			processTextContent += `<fo:table-cell>` + "\n"
			processTextContent += `<fo:block text-align="left">` + process.Description + `</fo:block>` + "\n"
			if cfg.Section("general").Key("showtranslations").String() == "yes" {
				for _, translation := range process.Translations {
					processTextContent += `<fo:block text-align="left" color="blue">` + translation + `</fo:block>` + "\n"
				}
			}
			processTextContent += `</fo:table-cell>` + "\n"
		case "translation":
			processTextContent += "<fo:table-cell>\n"
			if process.Translations != nil {
				for _, translation := range process.Translations {
					processTextContent += `<fo:block text-align="left">` + translation + `</fo:block>` + "\n"
				}
			} else {
				processTextContent += `<fo:block></fo:block>` + "\n"
			}
			processTextContent += "</fo:table-cell>\n"
		case "time":
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + process.Time + `</fo:block></fo:table-cell>` + "\n"
		case "tool":
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + process.Tool.Type + `</fo:block></fo:table-cell>` + "\n"
		case "torque":
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + process.Tool.Torque + `</fo:block></fo:table-cell>` + "\n"
		case "safety":
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + process.Tool.Class + `</fo:block></fo:table-cell>` + "\n"
		case "tolerance":
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + process.Tool.Tolerance + `</fo:block></fo:table-cell>` + "\n"
		case "socket":
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + process.Tool.Socket + `</fo:block></fo:table-cell>` + "\n"
		case "risk":
			processTextContent += `<fo:table-cell><fo:block text-align="left">` + process.Risk + `</fo:block></fo:table-cell>` + "\n"
		case "part":
			// format is name:number*quantity
			processTextContent += `<fo:table-cell>`
			if process.Parts != nil {
				for _, part := range process.Parts {
					partString := part.Number
					if part.Name != "" {
						if partString == "" {
							partString = part.Name
						} else {
							partString = part.Name + ":" + partString
						}
					}
					if part.Quantity != "" {
						partString += "*" + part.Quantity
					}
					processTextContent += `<fo:block text-align="center">` + partString + `</fo:block>`
				}
			} else {
				processTextContent += `<fo:block/>`
			}
			processTextContent += `</fo:table-cell>` + "\n"
		case "method":
			processTextContent += `<fo:table-cell><fo:block text-align="center">` + process.Check.Method + `</fo:block></fo:table-cell>` + "\n"
		case "criteria":
			processTextContent += `<fo:table-cell><fo:block text-align="left">` + process.Check.Criteria + `</fo:block></fo:table-cell>` + "\n"
		case "comment":
			if process.Comment.IsNoted == "yes" {
				backgroundColour = ` background-color="DeepSkyBlue"`
			} else {
				backgroundColour = ""
			}

			processTextContent += `<fo:table-cell` + backgroundColour + `>` + "\n"
			processTextContent += `<fo:block text-align="left">` + process.Comment.Text + `</fo:block>` + "\n"
			if cfg.Section("general").Key("showhcomment").String() == "yes" {
				// processTextContent += `<fo:block text-align="left" color="red"><fo:inline font-style="italic" font-weight="bold">` + process.Hcomment + `</fo:inline></fo:block>` + "\n"
				processTextContent += `<fo:block text-align="left" color="red">` + process.Hcomment + `</fo:block>` + "\n"
			}
			processTextContent += `</fo:table-cell>` + "\n"
		case "hcomment":
			processTextContent += `<fo:table-cell><fo:block text-align="left">` + process.Hcomment + `</fo:block></fo:table-cell>` + "\n"
		default:
			processTextContent += `<fo:table-cell><fo:block/></fo:table-cell>` + "\n"
		}
	}
	processTextContent += `</fo:table-row>` + "\n"

	content.ProcessTextContent = processTextContent
	processContent = append(processContent, content)

	if process.SubProcesses != nil {
		var subprocessNumber int = 0
		for _, subprocess := range process.SubProcesses {
			subprocessNumber++
			subprocessNumberString := processNumberString + "." + strconv.Itoa(subprocessNumber)
			subContent := processTableBodyContent(subprocess, subprocessNumberString)
			processContent = append(processContent, subContent...)
		}
	}

	return
}

func foXmlEnd() string {
	s := `
</fo:table>
<fo:block id="TheVeryLastPage" />
</fo:flow>
</fo:page-sequence>
</fo:root>
`

	return s
}
