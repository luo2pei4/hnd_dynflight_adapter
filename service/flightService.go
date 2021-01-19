package service

// Carrier 承运人
type Carrier struct {
	AirlineCD     string `json:"ＡＬコード"`
	AirlineJPName string `json:"ＡＬ和名称"`
	AirlineENName string `json:"ＡＬ英名称"`
	FlightNo      string `json:"便名"`
}

// RemarkReason 备考
type RemarkReason struct {
	RemarkJP   string `json:"ja"`
	RemarkEN   string `json:"en"`
	RemarkKO   string `json:"ko"`
	RemarkHans string `json:"zh-Hans"`
	RemarkHant string `json:"zh-Hant"`
}

// FlightInfo 航班信息
type FlightInfo struct {
	Carriers            []*Carrier    `json:"航空会社"`
	OrgnAirportCD       string        `json:"出発地空港コード"`
	OrgnAirportJPName   string        `json:"出発地空港和名称"`
	OrgnAirportENName   string        `json:"出発地空港英名称"`
	OrgnDirectionCD     string        `json:"出発地方面コード"`
	OrgnDirectionJPName string        `json:"出発地方面和名称"`
	OrgnDirectionENName string        `json:"出発地方面英名称"`
	DestAirportCD       string        `json:"行先地空港コード"`
	DestAirportJPName   string        `json:"行先地空港和名称"`
	DestAirportENName   string        `json:"行先地空港英名称"`
	DestDirectionCD     string        `json:"行先地方面コード"`
	DestDirectionJPName string        `json:"行先地方面和名称"`
	DestDirectionENName string        `json:"行先地方面英名称"`
	ViaAirportCD        string        `json:"経由地空港コード"`
	ViaAirportJPName    string        `json:"経由地空港和名称"`
	ViaAirportENName    string        `json:"経由地空港英名称"`
	ViaDirectionCD      string        `json:"経由地方面コード"`
	ViaDirectionJPName  string        `json:"経由地方面和名称"`
	ViaDirectionENName  string        `json:"経由地方面英名称"`
	ScheduleTime        string        `json:"定刻"`
	Status              string        `json:"状況"`
	ActualTile          string        `json:"変更時刻"`
	TerminalFlag        string        `json:"ターミナル区分"`
	SwingFlag           string        `json:"ウイング区分"`
	RemarkJPName        string        `json:"備考和名称"`
	RemarkENName        string        `json:"備考英名称"`
	Remark              *RemarkReason `json:"備考訳名称"`
	Fliker              string        `json:"フリッカ"`
	GateCD              string        `json:"ゲート番号コード"`
	GateJPName          string        `json:"ゲート和名称"`
	GateENName          string        `json:"ゲート英名称"`
	RemarkCD            string        `json:"備考コード"`
	CheckinCounter      string        `json:"チェックインカウンター番号"`
	SpotNo              string        `json:"スポット番号"`
	CraftType           string        `json:"機種コード"`
	OperatingStatus     string        `json:"運航状態"`
}

// FlightList 航班列表
type FlightList struct {
	LastUpdateTime string        `json:"last_upd"`
	List           []*FlightInfo `json:"flight_info"`
	FlightEnd      bool          `json:"flight_end"`
}
