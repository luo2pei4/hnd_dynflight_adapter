package service

import (
	"encoding/json"
	"fmt"
	"hda/dao"
	"hda/dto"
	"io/ioutil"
	"net/http"
	"time"
)

// TimeDiff 相差时间和时区
type TimeDiff struct {
	Hour int    `json:"text"`
	Gmt  string `json:"gmt"`
}

// Airport 机场信息
type Airport struct {
	Text  string   `json:"text"`
	Value string   `json:"value"`
	Tdiff TimeDiff `json:"timeDiff"`
}

// AirportItem 单项目
type AirportItem struct {
	Text     string     `json:"text"`
	Children []*Airport `json:"children"`
}

// AirportItems 项目列表
type AirportItems struct {
	Items []*AirportItem `json:"items"`
}

// AirportList 机场信息
type AirportList struct {
	JA     *AirportItems `json:"ja"`
	EN     *AirportItems `json:"en"`
	KO     *AirportItems `json:"ko"`
	ZhHans *AirportItems `json:"zh-Hans"`
	ZhHant *AirportItems `json:"zh-Hant"`
}

var airportMap map[string]*dto.AirportDto

// LoadAirports 加载机场信息
func LoadAirports() error {

	aMap, err := dao.QueryAirports()
	if err != nil {
		return err
	}

	airportMap = aMap

	fmt.Println("Load Airport info success.")

	return nil
}

// CrawlAirports 获取机场数据
func CrawlAirports() {

	timer := time.Tick(600 * 1e9)

	for {
		select {
		case <-timer:

			domAirportList, err := getAirportList("https://tokyo-haneda.com/site_resource/flight/data/dms/city_list_search.json")
			if err != nil {
				fmt.Println(err.Error())
			}

			aMap := make(map[string]*dto.AirportDto)

			if domAirportList != nil {
				editAirports(aMap, domAirportList)
			}

			intAirportList, err := getAirportList("https://tokyo-haneda.com/site_resource/flight/data/int/city_list_search.json")
			if err != nil {
				fmt.Println(err.Error())
			}

			if intAirportList != nil {
				editAirports(aMap, intAirportList)
			}

			failedCounter := 0

			for _, airport := range aMap {

				if airportMap[airport.IataCd] == nil {

					_, _, err := dao.SaveAirport(airport)

					if err != nil {
						failedCounter++
						fmt.Println(err.Error())
						continue
					}
				}
			}

			if failedCounter == 1 {
				fmt.Println("There is an company data insert failed")
			} else if failedCounter > 1 {
				fmt.Printf("There are %v company data insert failed", failedCounter)
			}
		}
	}
}

func getAirportList(url string) (airportList *AirportList, err error) {

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if response.Body != nil {

		body, _ := ioutil.ReadAll(response.Body)

		if body != nil {

			airportList = &AirportList{}
			err := json.Unmarshal(body, airportList)

			if err != nil {
				return nil, err
			}
		}
	}

	return
}

func editAirports(aMap map[string]*dto.AirportDto, airportList *AirportList) {

	items := airportList.JA.Items

	for _, item := range items {

		for _, info := range item.Children {
			if aMap[info.Value] == nil {
				aMap[info.Value] = &dto.AirportDto{
					NameJa: info.Text,
					IataCd: info.Value,
					Gmt:    info.Tdiff.Gmt,
				}
			} else {
				aMap[info.Value].NameJa = info.Text
			}
		}
	}

	items = airportList.EN.Items

	for _, item := range items {

		for _, info := range item.Children {
			if aMap[info.Value] == nil {
				aMap[info.Value] = &dto.AirportDto{
					NameEn: info.Text,
					IataCd: info.Value,
					Gmt:    info.Tdiff.Gmt,
				}
			} else {
				aMap[info.Value].NameEn = info.Text
			}
		}
	}

	items = airportList.KO.Items

	for _, item := range items {

		for _, info := range item.Children {
			if aMap[info.Value] == nil {
				aMap[info.Value] = &dto.AirportDto{
					NameKo: info.Text,
					IataCd: info.Value,
					Gmt:    info.Tdiff.Gmt,
				}
			} else {
				aMap[info.Value].NameKo = info.Value
			}
		}
	}

	items = airportList.ZhHans.Items

	for _, item := range items {

		for _, info := range item.Children {
			if aMap[info.Value] == nil {
				aMap[info.Value] = &dto.AirportDto{
					NameHans: info.Text,
					IataCd:   info.Value,
					Gmt:      info.Tdiff.Gmt,
				}
			} else {
				aMap[info.Value].NameHans = info.Text
			}
		}
	}

	items = airportList.ZhHant.Items

	for _, item := range items {

		for _, info := range item.Children {
			if aMap[info.Value] == nil {
				aMap[info.Value] = &dto.AirportDto{
					NameHant: info.Text,
					IataCd:   info.Value,
					Gmt:      info.Tdiff.Gmt,
				}
			} else {
				aMap[info.Value].NameHant = info.Text
			}
		}
	}
}
