package dao

import (
	"fmt"
	"hda/db"
	"hda/dto"
)

// SaveCampany 保存航空公司信息
func SaveCampany(dto *dto.CompanyDto) (lastInsertID, rowsAffected int64, err error) {

	sql := `INSERT INTO M_Company(name_en, name_ja, name_ko, name_hans, name_hant, icaocd, iatacd, homepage, countrycd) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')`
	sql = fmt.Sprintf(sql, dto.NameEn, dto.NameJa, dto.NameKo, dto.NameHans, dto.NameHant, dto.IcaoCd, dto.IataCd, dto.HomePage, dto.CountryCD)
	lastInsertID, rowsAffected, err = db.Execute(sql)

	return
}

// QueryCompanies 查询航空公司信息
func QueryCompanies() (companyMap map[string]*dto.CompanyDto, err error) {

	sql := `SELECT name_en, name_ja, name_ko, name_hans, name_hant, icaocd, iatacd, homepage, countrycd FROM adsb.M_Company`
	rows, err := db.Select(sql)

	if err != nil {
		return nil, err
	}

	companyMap = make(map[string]*dto.CompanyDto)

	for rows.Next() {
		dto := dto.CompanyDto{}
		rows.Scan(&dto.NameEn, &dto.NameJa, &dto.NameKo, &dto.NameHans, &dto.NameHant, &dto.IcaoCd, &dto.IataCd, &dto.HomePage, &dto.CountryCD)
		companyMap[dto.IcaoCd] = &dto
	}

	return
}
