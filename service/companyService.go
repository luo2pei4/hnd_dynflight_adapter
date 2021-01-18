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

// CompanyInfo 航空公司信息
type CompanyInfo struct {
	Name     string `json:"text"`
	IcaoCode string `json:"value"`
	IataCode string `json:"value02"`
	HomePage string `json:"url"`
}

// CompanyItems 航空公司列表
type CompanyItems struct {
	Items []*CompanyInfo `json:"items"`
}

// CompanyList 各种语言的航空公司列表
type CompanyList struct {
	JP     *CompanyItems `json:"ja"`
	EN     *CompanyItems `json:"en"`
	KO     *CompanyItems `json:"ko"`
	ZhHans *CompanyItems `json:"zh-Hans"`
	ZhHant *CompanyItems `json:"zh-Hant"`
}

var companyMap map[string]*dto.CompanyDto

// LoadCompanies 加载航空公司信息
func LoadCompanies() error {

	cMap, err := dao.QueryCompanies()

	if err != nil {
		return err
	}

	companyMap = cMap

	fmt.Println("Load Company info success.")

	return nil
}

// CrawlCompany 获取航空公司信息
func CrawlCompany() {

	timer := time.Tick(600 * 1e9)

	for {
		select {
		case <-timer:

			domCompanyList, err := getCompanyList("https://tokyo-haneda.com/site_resource/flight/data/dms/company_list_search.json")
			if err != nil {
				fmt.Println(err.Error())
			}

			cMap := make(map[string]*dto.CompanyDto)

			if domCompanyList != nil {
				editCompanies(cMap, domCompanyList)
			}

			intCompanyList, err := getCompanyList("https://tokyo-haneda.com/site_resource/flight/data/int/company_list_search.json")
			if err != nil {
				fmt.Println(err.Error())
			}

			if intCompanyList != nil {
				editCompanies(cMap, intCompanyList)
			}

			failedCounter := 0

			for _, company := range cMap {

				if companyMap[company.IcaoCd] == nil {

					_, _, err := dao.SaveCampany(company)

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

func getCompanyList(url string) (companyList *CompanyList, err error) {

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if response.Body != nil {

		body, _ := ioutil.ReadAll(response.Body)

		if body != nil {

			companyList = &CompanyList{}
			err := json.Unmarshal(body, companyList)

			if err != nil {
				return nil, err
			}
		}
	}

	return
}

func editCompanies(cMap map[string]*dto.CompanyDto, companyList *CompanyList) {

	items := companyList.JP.Items

	for _, info := range items {

		if cMap[info.IcaoCode] == nil {
			cMap[info.IcaoCode] = &dto.CompanyDto{
				NameJa:   info.Name,
				IcaoCd:   info.IcaoCode,
				IataCd:   info.IataCode,
				HomePage: info.HomePage,
			}
		} else {
			cMap[info.IcaoCode].NameJa = info.Name
		}
	}

	items = companyList.EN.Items

	for _, info := range items {

		if cMap[info.IcaoCode] == nil {
			cMap[info.IcaoCode] = &dto.CompanyDto{
				NameEn:   info.Name,
				IcaoCd:   info.IcaoCode,
				IataCd:   info.IataCode,
				HomePage: info.HomePage,
			}
		} else {
			cMap[info.IcaoCode].NameEn = info.Name
		}
	}

	items = companyList.KO.Items

	for _, info := range items {

		if cMap[info.IcaoCode] == nil {
			cMap[info.IcaoCode] = &dto.CompanyDto{
				NameKo:   info.Name,
				IcaoCd:   info.IcaoCode,
				IataCd:   info.IataCode,
				HomePage: info.HomePage,
			}
		} else {
			cMap[info.IcaoCode].NameKo = info.Name
		}
	}

	items = companyList.ZhHans.Items

	for _, info := range items {

		if cMap[info.IcaoCode] == nil {
			cMap[info.IcaoCode] = &dto.CompanyDto{
				NameHans: info.Name,
				IcaoCd:   info.IcaoCode,
				IataCd:   info.IataCode,
				HomePage: info.HomePage,
			}
		} else {
			cMap[info.IcaoCode].NameHans = info.Name
		}
	}

	items = companyList.ZhHant.Items

	for _, info := range items {

		if cMap[info.IcaoCode] == nil {
			cMap[info.IcaoCode] = &dto.CompanyDto{
				NameHant: info.Name,
				IcaoCd:   info.IcaoCode,
				IataCd:   info.IataCode,
				HomePage: info.HomePage,
			}
		} else {
			cMap[info.IcaoCode].NameHant = info.Name
		}
	}
}
