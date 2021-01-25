package service

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"hda/dao"
	"hda/dto"
	"io"
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

type statInfo struct {
	url            string
	arvDep         string
	lastUpdateTime string
	flightMap      map[string]string
}

var flightMap map[string]string

// hashValue 航班信息的哈希值
func (f *FlightInfo) hashValue() string {

	h := sha1.New()
	b, err := json.Marshal(*f)
	if err != nil {
		return ""
	}

	io.WriteString(h, string(b))

	return string(h.Sum(nil))
}

// CrawlHndDynFlight 爬取羽田机场航班动态
func CrawlHndDynFlight(cancel context.CancelFunc) {

	defer func() {

		if e := recover(); e != nil {
			fmt.Printf("Panicing %s\r\n", e)
			cancel()
		}
	}()

	timer := time.Tick(60 * 1e9)
	// dmsArvStat := &statInfo{
	// 	url:            "https://tokyo-haneda.com/app_resource/flight/data/dms/hdacfarv.json",
	// 	arvDep:         "ARV",
	// 	lastUpdateTime: "",
	// 	flightMap:      make(map[string]string, 1000),
	// }
	// dmsDepStat := &statInfo{
	// 	url:            "https://tokyo-haneda.com/app_resource/flight/data/dms/hdacfdep.json",
	// 	arvDep:         "DEP",
	// 	lastUpdateTime: "",
	// 	flightMap:      make(map[string]string, 1000),
	// }
	intArvStat := &statInfo{
		url:            "https://tokyo-haneda.com/app_resource/flight/data/int/hdacfarv.json",
		arvDep:         "ARV",
		lastUpdateTime: "",
		flightMap:      make(map[string]string, 1000),
	}
	intDepStat := &statInfo{
		url:            "https://tokyo-haneda.com/app_resource/flight/data/int/hdacfdep.json",
		arvDep:         "DEP",
		lastUpdateTime: "",
		flightMap:      make(map[string]string, 1000),
	}

	for {
		select {
		case <-timer:
			// hndDynFlight(dmsArvStat)
			// hndDynFlight(dmsDepStat)
			hndDynFlight(intArvStat)
			hndDynFlight(intDepStat)
		}
	}
}

// CrawlHndDynFlight 抓取羽田机场航班动态
func hndDynFlight(s *statInfo) {

	flightList, _, err := getHndDynFlightList(s.url)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	if s.lastUpdateTime == flightList.LastUpdateTime {
		return
	}

	fmt.Printf("%v before: %v, now: %v\n", s.url, s.lastUpdateTime, flightList.LastUpdateTime)

	// 保存最后更新时间
	s.lastUpdateTime = flightList.LastUpdateTime

	list := flightList.List

	dynFlightDto := &dto.HndDynFlightDto{}

	for _, info := range list {

		airlineCD := info.Carriers[0].AirlineCD
		flightNo := info.Carriers[0].FlightNo
		scheduleTime := info.ScheduleTime
		scheduleTime = strings.ReplaceAll(scheduleTime, "/", "")
		scheduleTime = strings.ReplaceAll(scheduleTime, " ", "")
		scheduleTime = strings.ReplaceAll(scheduleTime, ":", "")
		fuid := airlineCD + flightNo + "-" + scheduleTime + "-" + s.arvDep
		sha1Value := info.hashValue()

		if flightMap[fuid] == sha1Value {
			continue
		}

		var adminFlightID int64

		for idx, carrier := range info.Carriers {

			// administrating flight
			if idx == 0 {
				dynFlightDto.AirlineCD = carrier.AirlineCD
				dynFlightDto.FlightNo = carrier.FlightNo
				dynFlightDto.OrgnAirportCD = info.OrgnAirportCD
				dynFlightDto.OrgnDirectionCD = info.OrgnDirectionCD
				dynFlightDto.OrgnDirectionJPName = info.OrgnDirectionJPName
				dynFlightDto.OrgnDirectionENName = info.OrgnDirectionENName
				dynFlightDto.DestAirportCD = info.DestAirportCD
				dynFlightDto.DestDirectionCD = info.DestDirectionCD
				dynFlightDto.DestDirectionJPName = info.DestDirectionJPName
				dynFlightDto.DestDirectionENName = info.DestDirectionENName
				dynFlightDto.ViaAirportCD = info.ViaAirportCD
				dynFlightDto.ViaDirectionCD = info.ViaDirectionCD
				dynFlightDto.ViaDirectionJPName = info.ViaDirectionJPName
				dynFlightDto.ViaDirectionENName = info.ViaDirectionENName
				dynFlightDto.ScheduleTime = info.ScheduleTime
				dynFlightDto.ActualTime = info.ActualTime
				dynFlightDto.Status = info.Status
				dynFlightDto.Terminal = info.TerminalFlag
				dynFlightDto.Swing = info.SwingFlag
				dynFlightDto.RemarkJPName = info.RemarkJPName
				dynFlightDto.RemarkENName = info.RemarkENName
				dynFlightDto.RemarkJP = info.Remark.RemarkJP
				dynFlightDto.RemarkEN = info.Remark.RemarkEN
				dynFlightDto.RemarkKO = info.Remark.RemarkKO
				dynFlightDto.RemarkHans = info.Remark.RemarkHans
				dynFlightDto.RemarkHant = info.Remark.RemarkHant
				dynFlightDto.Fliker = info.Fliker
				dynFlightDto.GateCD = info.GateCD
				dynFlightDto.RemarkCD = info.RemarkCD
				dynFlightDto.CheckinCounter = info.CheckinCounter
				dynFlightDto.SpotNo = info.SpotNo
				dynFlightDto.CraftType = info.CraftType
				dynFlightDto.OperatingStatus = info.OperatingStatus
				dynFlightDto.Createtime = getCurrentTime()
				lastInsertID, _, err := dao.SaveHndDynFlight(dynFlightDto)

				if err != nil {
					fmt.Println("Save flight error:", err.Error())
					continue
				}

				_, err = dao.DeleteShareCode(lastInsertID)
				adminFlightID = lastInsertID

			} else {

				shareCodeInfo := &dto.HndShareCodeDto{
					AdminFlightID: adminFlightID,
					AirlineCD:     carrier.AirlineCD,
					FlightNo:      carrier.FlightNo,
				}
				_, _, err = dao.SaveShareCode(shareCodeInfo)

				if err != nil {
					fmt.Println("Save sharecode error:", err.Error())
					continue
				}
			}
		}

		flightMap[fuid] = sha1Value
	}

	fmt.Printf("%v , flightMap len: %v\n", s.url, len(flightMap))
}

func getHndDynFlightList(url string) (flightList *FlightList, msg string, err error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, "", err
	}

	// 设置连接关闭标志
	req.Close = true
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, "", err
	}

	if response.Body != nil {

		body, _ := ioutil.ReadAll(response.Body)

		if body != nil {

			msg = string(body)
			flightList = &FlightList{}
			err := json.Unmarshal(body, flightList)

			if err != nil {
				return nil, msg, err
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
