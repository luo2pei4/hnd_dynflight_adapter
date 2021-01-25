package dto

// HndDynFlightDto 羽田机场动态航班信息
type HndDynFlightDto struct {
	ID                  int
	AirlineCD           string
	FlightNo            string
	OrgnAirportCD       string
	OrgnDirectionCD     string
	OrgnDirectionJPName string
	OrgnDirectionENName string
	DestAirportCD       string
	DestDirectionCD     string
	DestDirectionJPName string
	DestDirectionENName string
	ViaAirportCD        string
	ViaDirectionCD      string
	ViaDirectionJPName  string
	ViaDirectionENName  string
	ScheduleTime        string
	ActualTime          string
	Status              string
	Terminal            string
	Swing               string
	RemarkJPName        string
	RemarkENName        string
	RemarkJP            string
	RemarkEN            string
	RemarkKO            string
	RemarkHans          string
	RemarkHant          string
	Fliker              string
	GateCD              string
	RemarkCD            string
	CheckinCounter      string
	SpotNo              string
	CraftType           string
	OperatingStatus     string
	Createtime          string
}

// HndShareCodeDto sharecode flight info
type HndShareCodeDto struct {
	ID            int
	AdminFlightID int64
	AirlineCD     string
	FlightNo      string
}
