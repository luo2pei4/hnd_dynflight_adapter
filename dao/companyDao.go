package dao

import (
	"fmt"
	"hda/dto"
)

// SaveCampany 保存航空公司信息
func SaveCampany(dto *dto.CompanyDto) (lastInsertID, rowsAffected int64, err error) {

	sql := `INSERT INTO M_Company(name_en, name_ja, name_ko, name_hans, name_hant, icaocd, iatacd, homepage, countrycd) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')`
	sql = fmt.Sprintf(sql, dto.NameEn, dto.NameJa, dto.NameKo, dto.NameHans, dto.NameHant, dto.IcaoCd, dto.IataCd, dto.HomePage, dto.CountryCD)
	lastInsertID, rowsAffected, err = conn.Insert(sql)

	return
}
