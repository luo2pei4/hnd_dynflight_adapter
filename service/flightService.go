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
	"strconv"
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
}

var flightMap map[string][]string

func init() {

	flightMap = make(map[string][]string, 4000)
	dmsArvStat := &statInfo{url: "https://tokyo-haneda.com/app_resource/flight/data/dms/hdacfarv.json", arvDep: "ARV"}
	dmsDepStat := &statInfo{url: "https://tokyo-haneda.com/app_resource/flight/data/dms/hdacfdep.json", arvDep: "DEP"}
	intArvStat := &statInfo{url: "https://tokyo-haneda.com/app_resource/flight/data/int/hdacfarv.json", arvDep: "ARV"}
	intDepStat := &statInfo{url: "https://tokyo-haneda.com/app_resource/flight/data/int/hdacfdep.json", arvDep: "DEP"}

	initializeFlightMap(dmsArvStat)
	initializeFlightMap(dmsDepStat)
	initializeFlightMap(intArvStat)
	initializeFlightMap(intDepStat)
}

func initializeFlightMap(s *statInfo) {

	fmt.Println(s.url)
	flightList, _, err := getHndFlightList(s.url)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	list := flightList.List
	for _, info := range list {

		carrierCD := info.Carriers[0].AirlineCD
		flightNo := info.Carriers[0].FlightNo
		scheduleTime := info.ScheduleTime
		scheduleTime = strings.ReplaceAll(scheduleTime, "/", "")
		scheduleTime = strings.ReplaceAll(scheduleTime, " ", "")
		scheduleTime = strings.ReplaceAll(scheduleTime, ":", "")
		var orgnAirportCD string
		var destAirportCD string
		var arvDepAirport string

		if s.arvDep == "ARV" {
			orgnAirportCD = info.OrgnAirportCD
			destAirportCD = "HND"
			arvDepAirport = info.OrgnAirportCD
		} else {
			orgnAirportCD = "HND"
			destAirportCD = info.DestAirportCD
			arvDepAirport = info.DestAirportCD
		}
		fuid := carrierCD + flightNo + "-" + scheduleTime + "-" + s.arvDep + "-" + arvDepAirport
		sha1Value := info.hashValue()

		fmt.Printf("FUID: %v, ", fuid)

		flightDto, err := dao.FindHndFlight(carrierCD, flightNo, orgnAirportCD, destAirportCD, scheduleTime)

		if err != nil {
			fmt.Println("Find flightId error,", err.Error())
			continue
		}

		fmt.Printf("FlightID: %v\n", flightDto.ID)
		flightMap[fuid] = []string{sha1Value, strconv.FormatInt(int64(flightDto.ID), 10), scheduleTime}
	}
}

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
	dmsArvStat := &statInfo{url: "https://tokyo-haneda.com/app_resource/flight/data/dms/hdacfarv.json", arvDep: "ARV"}
	dmsDepStat := &statInfo{url: "https://tokyo-haneda.com/app_resource/flight/data/dms/hdacfdep.json", arvDep: "DEP"}
	intArvStat := &statInfo{url: "https://tokyo-haneda.com/app_resource/flight/data/int/hdacfarv.json", arvDep: "ARV"}
	intDepStat := &statInfo{url: "https://tokyo-haneda.com/app_resource/flight/data/int/hdacfdep.json", arvDep: "DEP"}

	for {
		select {
		case <-timer:
			hndDynFlight(dmsArvStat)
			hndDynFlight(dmsDepStat)
			hndDynFlight(intArvStat)
			hndDynFlight(intDepStat)
			cleanFlightMap()
		}
	}
}

// CrawlHndDynFlight 抓取羽田机场航班动态
func hndDynFlight(s *statInfo) {

	beginTime := time.Now().UnixNano()
	fmt.Printf("%v\nbegin, %v\n", s.url, time.Unix(0, beginTime))

	flightList, _, err := getHndFlightList(s.url)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	if s.lastUpdateTime == flightList.LastUpdateTime {
		fmt.Printf("end.., %v\n", time.Unix(0, time.Now().UnixNano()))
		return
	}

	// 保存最后更新时间
	s.lastUpdateTime = flightList.LastUpdateTime
	list := flightList.List

	// dynFlightDto := &dto.HndDynFlightDto{}
	hndFlightDto := &dto.HndFlightDto{}

	for _, info := range list {

		carrierCD := info.Carriers[0].AirlineCD
		flightNo := info.Carriers[0].FlightNo
		scheduleTime := info.ScheduleTime
		scheduleTime = strings.ReplaceAll(scheduleTime, "/", "")
		scheduleTime = strings.ReplaceAll(scheduleTime, " ", "")
		scheduleTime = strings.ReplaceAll(scheduleTime, ":", "")
		var arvDepAirport string
		if s.arvDep == "DEP" {
			arvDepAirport = info.DestAirportCD
		} else {
			arvDepAirport = info.OrgnAirportCD
		}
		fuid := carrierCD + flightNo + "-" + scheduleTime + "-" + s.arvDep + "-" + arvDepAirport
		sha1Value := info.hashValue()

		if flightMap[fuid] == nil || len(flightMap[fuid]) == 0 {
			flightMap[fuid] = make([]string, 3)
		} else {
			if flightMap[fuid][0] == sha1Value {
				continue
			}
		}

		var adminFlightID int64
		hasShareCode := false

		if len(info.Carriers) > 1 {
			hasShareCode = true
		}

		for idx, carrier := range info.Carriers {

			// administrating flight
			if idx == 0 {

				hndFlightDto.CarrierCD = carrier.AirlineCD
				hndFlightDto.FlightNo = carrier.FlightNo
				hndFlightDto.CraftType = info.CraftType

				if s.arvDep == "ARV" {
					hndFlightDto.OrgnAirportCd = info.OrgnAirportCD
					hndFlightDto.DestAirportCd = "HND"
				} else {
					hndFlightDto.OrgnAirportCd = "HND"
					hndFlightDto.DestAirportCd = info.DestAirportCD
				}

				hndFlightDto.ViaAirportCd = info.ViaAirportCD
				hndFlightDto.ScheduleTime = info.ScheduleTime
				flightMap[fuid][2] = info.ScheduleTime

				hndFlightDto.ActualTime = info.ActualTime
				hndFlightDto.Terminal = info.TerminalFlag
				hndFlightDto.Swing = info.SwingFlag

				remarkEn := info.Remark.RemarkEN
				if strings.Contains(remarkEn, "EXIT") {
					hndFlightDto.ExitCD = remarkEn[4:]
				} else {
					hndFlightDto.Status = remarkEn
				}

				hndFlightDto.GateCd = info.GateCD
				hndFlightDto.CheckinCounter = info.CheckinCounter
				hndFlightDto.SpotNo = info.SpotNo
				hndFlightDto.CreateTime = getCurrentTime()

				if flightMap[fuid][1] == "" {

					lastInsertID, _, err := dao.SaveHndFlight(hndFlightDto)

					if err != nil {
						fmt.Println("Save flight error:", err.Error())
						break
					}

					adminFlightID = lastInsertID
					flightMap[fuid][1] = strconv.FormatInt(lastInsertID, 10)

				} else {

					adminFlightID, e := strconv.ParseInt(flightMap[fuid][1], 10, 64)

					if e != nil {
						fmt.Printf("FlightID: %v convert failed.\n", flightMap[fuid][1])
						break
					}

					hndFlightDto.ID = adminFlightID
					_, _, err = dao.SaveHndFlightChanges(adminFlightID)
					if err != nil {
						fmt.Println("Save flight changes error:", err.Error())
						break
					}
					fmt.Printf("%v flight changes save OK, ", adminFlightID)

					_, err = dao.UpdateHndFlight(hndFlightDto)
					if err != nil {
						fmt.Println("Update flight error:", err.Error())
						break
					}
					fmt.Printf("flight info update OK, ")

					if hasShareCode {
						rowsAffected, err := dao.DeleteShareCode(adminFlightID)
						if err != nil {
							fmt.Println("Delete sharecode error:", err.Error())
							break
						}
						fmt.Printf("flight sharecode %v record(s) delete OK.\n", rowsAffected)
					}
				}

			} else {

				shareCodeInfo := &dto.HndShareCodeDto{
					AdminFlightID: adminFlightID,
					AirlineCD:     carrier.AirlineCD,
					FlightNo:      carrier.FlightNo,
				}
				lastInsertID, _, err := dao.SaveShareCode(shareCodeInfo)

				if err != nil {
					fmt.Println("Save sharecode error:", err.Error())
					break
				}
				fmt.Printf("%v flight sharecode save OK, ID: %v\n", adminFlightID, lastInsertID)
			}
		}

		flightMap[fuid][0] = sha1Value
	}

	endTime := time.Now().UnixNano()
	diff := endTime - beginTime
	cost := float64(diff) / 1e6

	fmt.Printf("end.., %v, cost:%v\n", time.Unix(0, endTime), cost)
}

func getHndFlightList(url string) (flightList *FlightList, msg string, err error) {

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

func cleanFlightMap() {

	fmt.Printf("Clean flightMap, before: %v, ", len(flightMap))

	temp := time.Now().Add(-48 * time.Hour)
	keyList := make([]string, 0)

	for key, s := range flightMap {

		strTime := s[2]

		scheduleTime, err := time.Parse("2006/01/02 15:04:05", strTime)
		if err != nil {
			fmt.Println("Schedule time parse error.", err.Error())
			continue
		}

		if scheduleTime.Before(temp) {
			keyList = append(keyList, key)
		}
	}

	if len(keyList) > 0 {

		for _, key := range keyList {
			delete(flightMap, key)
		}
	}

	fmt.Printf("after: %v\n", len(flightMap))
}
