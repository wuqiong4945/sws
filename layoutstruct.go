package main

type AreaStruct struct {
	Name     string
	Paper    PaperStruct
	Position PositionStruct
	Stations []StationStruct
}

type PaperStruct struct {
	W, H int
}

type PositionStruct struct {
	X, Y   int
	R      float64
	W, H   int
	VW, VH int
	Kind   string
}

type StationStruct struct {
	Name     string
	Position PositionStruct
	Swses    []SwsStruct
}

type OperatorInfoStruct struct {
	Position      string
	OperationTime OperationTimeStruct
	// SwsContent    *OperatorStruct
}
