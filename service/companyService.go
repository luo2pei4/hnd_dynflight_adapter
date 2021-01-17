package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"hda/dao"
	"hda/dto"
	"io/ioutil"
	"net/http"
)

// CompanyInfo 航空公司信息
type CompanyInfo struct {
	Name     string `json:"text"`
	IcaoCode string `json:"value"`
	IataCode string `json:"value02"`
	HomePage string `json:"url"`
}

// Items 航空公司列表
type Items struct {
	List []*CompanyInfo `json:"items"`
}

// CompanyList 各种语言的航空公司列表
type CompanyList struct {
	JP     *Items `json:"ja"`
	EN     *Items `json:"en"`
	KO     *Items `json:"ko"`
	ZhHans *Items `json:"zh-Hans"`
	ZhHant *Items `json:"zh-Hant"`
}

// CrawlCompany 获取航空公司信息
func CrawlCompany() error {

	domCompanyList, err := crawl("https://tokyo-haneda.com/site_resource/flight/data/dms/company_list_search.json")
	if err != nil {
		return err
	}

	cMap := make(map[string]*dto.CompanyDto)
	edit(cMap, domCompanyList)

	intCompanyList, err := crawl("https://tokyo-haneda.com/site_resource/flight/data/int/company_list_search.json")
	if err != nil {
		return err
	}

	edit(cMap, intCompanyList)

	failedCounter := 0

	for _, company := range cMap {

		_, _, err := dao.SaveCampany(company)
		if err != nil {
			failedCounter++
			fmt.Println(err.Error())
			continue
		}
	}

	if failedCounter == 1 {
		return errors.New("There is an company data insert failed")
	} else if failedCounter > 1 {
		return fmt.Errorf("There are %v company data insert failed", failedCounter)
	}

	return nil
}

func crawl(url string) (companyList *CompanyList, err error) {

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

func edit(cMap map[string]*dto.CompanyDto, companyList *CompanyList) {

	list := companyList.JP.List

	for _, info := range list {

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

	list = companyList.EN.List

	for _, info := range list {

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

	list = companyList.KO.List

	for _, info := range list {

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

	list = companyList.ZhHans.List

	for _, info := range list {

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

	list = companyList.ZhHant.List

	for _, info := range list {

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
