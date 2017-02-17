package main

type SwsStruct struct {
	Info     InfoStruct     `xml:"info"`
	Operator OperatorStruct `xml:"operator"`
}

type InfoStruct struct {
	Author         string `xml:"author,attr"`
	Department     string `xml:"department,attr"`
	UpdateTime     string `xml:"updatetime,attr"`
	Column         string `xml:"column,attr"`
	Title          string `xml:"title,attr"`
	AdditionalInfo string `xml:"additionalinfo,attr"`
}

type OperatorStruct struct {
	Model    string  `xml:"model,attr"`
	Station  string  `xml:"station,attr"`
	Position string  `xml:"position,attr"`
	Wtime    float32 `xml:"wtime,attr"`

	FirstProcessNumber int `xml:"firstProcessNumber,attr"`

	Safety    SafetyStruct    `xml:"safety"`
	Processes []ProcessStruct `xml:"process"`
}

type SafetyStruct struct {
	IsESDShoes      string `xml:"isESDShoes,attr"`
	IsWorkware      string `xml:"isWorkware,attr"`
	IsSafetyGlasses string `xml:"isSafetyGlasses,attr"`
	IsSafetyGloves  string `xml:"isSafetyGloves,attr"`
}

type ProcessStruct struct {
	Image     string  `xml:"image,attr"`
	ImageSize string  `xml:"imagesize,attr"`
	Option    string  `xml:"option,attr"`
	Tvg       string  `xml:"tvg,attr"`
	IsKey     string  `xml:"isKey,attr"`
	Time      float32 `xml:"time,attr"`
	Nvtime    float32 `xml:"nvtime,attr"`

	Description  string        `xml:"description"`
	Translations []string      `xml:"translation"`
	Parts        []PartStruct  `xml:"part"`
	Tool         ToolStruct    `xml:"tool"`
	Risk         string        `xml:"risk"`
	Check        CheckStruct   `xml:"check"`
	Comment      CommentStruct `xml:"comment"`
	Hcomment     string        `xml:"hcomment"`

	SubProcesses []ProcessStruct `xml:"subprocess"`
}

type PartStruct struct {
	Number   string `xml:"number,attr"`
	Quantity string `xml:"quantity,attr"`
	Family   string `xml:"family,attr"`
	Name     string `xml:"name,attr"`
}

type ToolStruct struct {
	Type      string `xml:"type,attr"`
	Torque    string `xml:"torque,attr"`
	Class     string `xml:"class,attr"`
	Tolerance string `xml:"tolerance,attr"`
	Socket    string `xml:"socket,attr"`
}

type CheckStruct struct {
	Method   string `xml:"method"`
	Criteria string `xml:"criteria"`
}

type CommentStruct struct {
	IsNoted string `xml:"isNoted,attr"`
	Text    string `xml:",chardata"`
}
