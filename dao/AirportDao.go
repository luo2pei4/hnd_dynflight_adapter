package dao

import (
	"fmt"
	"hda/dto"
)

// SaveAirport 保存机场信息
func SaveAirport(dto *dto.AirportDto) (lastInsertID, rowsAffected int64, err error) {

	sql := `INSERT INTO M_Airport(name_en, name_ja, name_ko, name_hans, name_hant, icaocd, iatacd, gmt, region) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')`
	sql = fmt.Sprintf(sql, dto.NameEn, dto.NameJa, dto.NameKo, dto.NameHans, dto.NameHant, dto.IcaoCd, dto.IataCd, dto.Gmt, dto.Region)
	lastInsertID, rowsAffected, err = conn.Insert(sql)

	return
}

// QueryAirports 查询机场信息
func QueryAirports() (airportMap map[string]*dto.AirportDto, err error) {

	sql := `SELECT name_en, name_ja, name_ko, name_hans, name_hant, icaocd, iatacd, gmt, region	FROM adsb.M_Airport`
	rows, err := conn.Select(sql)

	if err != nil {
		return nil, err
	}

	airportMap = make(map[string]*dto.AirportDto)

	for rows.Next() {
		dto := dto.AirportDto{}
		rows.Scan(&dto.NameEn, &dto.NameJa, &dto.NameKo, &dto.NameHans, &dto.NameHant, &dto.IcaoCd, &dto.IataCd, &dto.Gmt, &dto.Region)
		airportMap[dto.IataCd] = &dto
	}

	return
}
