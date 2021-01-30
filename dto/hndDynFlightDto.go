package dto

// HndFlightDto 羽田航班信息
type HndFlightDto struct {
	ID             int64
	CarrierCD      string
	FlightNo       string
	CraftType      string
	OrgnAirportCd  string
	DestAirportCd  string
	ViaAirportCd   string
	ScheduleTime   string
	ActualTime     string
	Terminal       string
	Swing          string
	Status         string
	GateCd         string
	CheckinCounter string
	ExitCD         string
	SpotNo         string
	CreateTime     string
}

// HndShareCodeDto sharecode flight info
type HndShareCodeDto struct {
	ID            int
	AdminFlightID int64
	AirlineCD     string
	FlightNo      string
}
