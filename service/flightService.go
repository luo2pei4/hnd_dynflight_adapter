package service

import (
	"encoding/json"
	"fmt"
	"hda/dao"
	"hda/dto"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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
	ActualTime          string        `json:"変更時刻"`
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

// CrawlHndDynFlight 抓取羽田机场航班动态
func CrawlHndDynFlight(url string) {

	flightList, err := getHndDynFlightList(url)

	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	list := flightList.List

	dto := &dto.HndDynFlightDto{}

	for _, info := range list {

		for _, carrier := range info.Carriers {
			dto.AirlineCD = carrier.AirlineCD
			dto.FlightNo = carrier.FlightNo
			dto.OrgnAirportCD = info.OrgnAirportCD
			dto.OrgnDirectionCD = info.OrgnDirectionCD
			dto.OrgnDirectionJPName = info.OrgnDirectionJPName
			dto.OrgnDirectionENName = info.OrgnDirectionENName
			dto.DestAirportCD = info.DestAirportCD
			dto.DestDirectionCD = info.DestDirectionCD
			dto.DestDirectionJPName = info.DestDirectionJPName
			dto.DestDirectionENName = info.DestDirectionENName
			dto.ViaAirportCD = info.ViaAirportCD
			dto.ViaDirectionCD = info.ViaDirectionCD
			dto.ViaDirectionJPName = info.ViaDirectionJPName
			dto.ViaDirectionENName = info.ViaDirectionENName
			dto.ScheduleTime = info.ScheduleTime
			dto.ActualTime = info.ActualTime
			dto.Status = info.Status
			dto.Terminal = info.TerminalFlag
			dto.Swing = info.SwingFlag
			dto.RemarkJPName = info.RemarkJPName
			dto.RemarkENName = info.RemarkENName
			dto.RemarkJP = info.Remark.RemarkJP
			dto.RemarkEN = info.Remark.RemarkEN
			dto.RemarkKO = info.Remark.RemarkKO
			dto.RemarkHans = info.Remark.RemarkHans
			dto.RemarkHant = info.Remark.RemarkHant
			dto.Fliker = info.Fliker
			dto.GateCD = info.GateCD
			dto.RemarkCD = info.RemarkCD
			dto.CheckinCounter = info.CheckinCounter
			dto.SpotNo = info.SpotNo
			dto.CraftType = info.CraftType
			dto.OperatingStatus = info.OperatingStatus
			dto.Createtime = getCurrentTime()
			_, _, err := dao.SaveHndDynFlight(dto)

			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func getHndDynFlightList(url string) (flightList *FlightList, err error) {

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if response.Body != nil {

		body, _ := ioutil.ReadAll(response.Body)

		if body != nil {

			flightList = &FlightList{}
			err := json.Unmarshal(body, flightList)

			if err != nil {
				return nil, err
			}
		}
	}

	return
}

func getCurrentTime() string {

	now := time.Now().Local().Format(time.RFC3339)
	currentTime := now[:19]
	currentTime = strings.ReplaceAll(currentTime, "T", " ")

	return currentTime
}
