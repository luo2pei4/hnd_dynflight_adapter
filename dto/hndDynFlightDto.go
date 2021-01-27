package dto

// HndDynFlightDto 羽田机场动态航班信息
// type HndDynFlightDto struct {
// 	ID                  int
// 	AirlineCD           string
// 	FlightNo            string
// 	OrgnAirportCD       string
// 	OrgnDirectionCD     string
// 	OrgnDirectionJPName string
// 	OrgnDirectionENName string
// 	DestAirportCD       string
// 	DestDirectionCD     string
// 	DestDirectionJPName string
// 	DestDirectionENName string
// 	ViaAirportCD        string
// 	ViaDirectionCD      string
// 	ViaDirectionJPName  string
// 	ViaDirectionENName  string
// 	ScheduleTime        string
// 	ActualTime          string
// 	Status              string
// 	Terminal            string
// 	Swing               string
// 	RemarkJPName        string
// 	RemarkENName        string
// 	RemarkJP            string
// 	RemarkEN            string
// 	RemarkKO            string
// 	RemarkHans          string
// 	RemarkHant          string
// 	Fliker              string
// 	GateCD              string
// 	RemarkCD            string
// 	CheckinCounter      string
// 	SpotNo              string
// 	CraftType           string
// 	OperatingStatus     string
// 	Createtime          string
// }

// HndFlightDto 羽田航班信息
type HndFlightDto struct {
	ID             int
	CarrierCd      string
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
	ExitCd         string
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
