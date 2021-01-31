package dao

import (
	"fmt"
	"hda/db"
	"hda/dto"
	"time"
)

// SaveHndFlight 保存羽田机场航班
func SaveHndFlight(hndFlightDto *dto.HndFlightDto) (lastInsertID, rowsAffected int64, err error) {

	sql := `INSERT INTO adsb.hnd_flight (
			carriercd, flightno, crafttype, 
			orgnairportcd, destairportcd, viaairportcd, 
			scheduletime, actualtime, 
			terminal, swing, status, gatecd, checkincounter, exitcd, spotno, createtime, lastupdtime
		) VALUES (
			'%s', '%s', '%s', 
			'%s', '%s', '%s', 
			%s, %s, 
			'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'
		)`

	scheduleTime := hndFlightDto.ScheduleTime
	actualTime := hndFlightDto.ActualTime

	if scheduleTime == "" {
		scheduleTime = "null"
	} else {
		scheduleTime = "'" + scheduleTime + "'"
	}

	if actualTime == "" {
		actualTime = "null"
	} else {
		actualTime = "'" + actualTime + "'"
	}

	sql = fmt.Sprintf(sql,
		hndFlightDto.CarrierCD, hndFlightDto.FlightNo, hndFlightDto.CraftType,
		hndFlightDto.OrgnAirportCd, hndFlightDto.DestAirportCd, hndFlightDto.ViaAirportCd,
		scheduleTime, actualTime,
		hndFlightDto.Terminal, hndFlightDto.Swing, hndFlightDto.Status, hndFlightDto.GateCd, hndFlightDto.CheckinCounter, hndFlightDto.ExitCD, hndFlightDto.SpotNo,
		hndFlightDto.CreateTime, hndFlightDto.CreateTime,
	)

	lastInsertID, rowsAffected, err = db.Execute(sql)

	return
}

// SaveShareCode save sharecode flight info
func SaveShareCode(shareCode *dto.HndShareCodeDto) (lastInsertID, rowsAffected int64, err error) {

	sql := `INSERT INTO adsb.hnd_flight_sharecode (adminflightid, airlinecd, flightno) VALUES (%v, '%v', '%v')`
	sql = fmt.Sprintf(sql, shareCode.AdminFlightID, shareCode.AirlineCD, shareCode.FlightNo)
	lastInsertID, rowsAffected, err = db.Execute(sql)

	return
}

// DeleteShareCode delete sharecode info
func DeleteShareCode(adminFlightID int64) (rowsAffected int64, err error) {

	sql := `SELECT count(1) as counter FROM adsb.hnd_flight_sharecode where adminflightid = %v`
	sql = fmt.Sprintf(sql, adminFlightID)
	rows, err := db.Select(sql)

	if err != nil {
		return 0, err
	}

	counter := 0
	rows.Next()
	rows.Scan(&counter)
	rows.Close()

	if counter == 0 {
		return 0, nil
	}

	sql = `DELETE FROM adsb.hnd_flight_sharecode WHERE adminflightid = %v`
	sql = fmt.Sprintf(sql, adminFlightID)
	_, rowsAffected, err = db.Execute(sql)

	return
}

// UpdateHndFlight 更新羽田航班信息
func UpdateHndFlight(hndFlightDto *dto.HndFlightDto) (rowsAffected int64, err error) {
	sql := `UPDATE 
				adsb.hnd_flight
			SET 
				carriercd='%s', flightno='%s', crafttype='%s', 
				orgnairportcd='%s', destairportcd='%s', viaairportcd='%s', 
				scheduletime=%s, actualtime=%s, 
				terminal='%s', swing='%s', status='%s', gatecd='%s', checkincounter='%s', exitcd='%s', spotno='%s', 
				lastupdtime='%s'
			WHERE 
				id=%v`

	scheduleTime := hndFlightDto.ScheduleTime
	actualTime := hndFlightDto.ActualTime

	if scheduleTime == "" {
		scheduleTime = "null"
	} else {
		scheduleTime = "'" + scheduleTime + "'"
	}

	if actualTime == "" {
		actualTime = "null"
	} else {
		actualTime = "'" + actualTime + "'"
	}

	lastupdtime := time.Now().Format("2006-01-02 15:04:05")

	sql = fmt.Sprintf(sql,
		hndFlightDto.CarrierCD, hndFlightDto.FlightNo, hndFlightDto.CraftType,
		hndFlightDto.OrgnAirportCd, hndFlightDto.DestAirportCd, hndFlightDto.ViaAirportCd,
		scheduleTime, actualTime,
		hndFlightDto.Terminal, hndFlightDto.Swing, hndFlightDto.Status, hndFlightDto.GateCd, hndFlightDto.CheckinCounter, hndFlightDto.ExitCD, hndFlightDto.SpotNo,
		lastupdtime,
		hndFlightDto.ID,
	)

	_, rowsAffected, err = db.Execute(sql)

	return
}

// FindHndFlight 根据航班号, 始发目的站和计划时间查询航班
func FindHndFlight(carrierCd, flightNo, orgnAirportCd, destAirportCD, scheduleTime string) (hndFlightDto *dto.HndFlightDto, err error) {
	sql := `SELECT 
				id, carriercd, flightno, crafttype, 
				orgnairportcd, destairportcd, viaairportcd, 
				scheduletime, case when actualtime is null then '' else actualtime end as actualtime, 
				terminal, swing, status, gatecd, checkincounter, exitcd, spotno 
			FROM 
				adsb.hnd_flight
			WHERE 
			carriercd = '%s'
			AND flightno = '%s'
			AND orgnairportcd = '%s'
			AND destAirportCD = '%s'
			AND scheduletime = '%s'`
	sql = fmt.Sprintf(sql, carrierCd, flightNo, orgnAirportCd, destAirportCD, scheduleTime)

	rows, err := db.Select(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	hndFlightDto = &dto.HndFlightDto{}
	rows.Next()
	err = rows.Scan(&hndFlightDto.ID, &hndFlightDto.CarrierCD, &hndFlightDto.FlightNo, &hndFlightDto.CraftType,
		&hndFlightDto.OrgnAirportCd, &hndFlightDto.DestAirportCd, &hndFlightDto.ViaAirportCd,
		&hndFlightDto.ScheduleTime, &hndFlightDto.ActualTime,
		&hndFlightDto.Terminal, &hndFlightDto.Swing, &hndFlightDto.Status, &hndFlightDto.GateCd, &hndFlightDto.CheckinCounter, &hndFlightDto.ExitCD, &hndFlightDto.SpotNo,
	)
	if err != nil {
		return nil, err
	}

	return
}

// SaveHndFlightChanges 保存上一次变更后的航班数据
func SaveHndFlightChanges(adminFlightID int64) (lastInsertID, rowsAffected int64, err error) {
	sql := `INSERT INTO adsb.hnd_flight_changes (
		adminflightid, carriercd, flightno, crafttype, orgnairportcd, destairportcd, viaairportcd, scheduletime, actualtime, terminal, swing, status, gatecd, checkincounter, exitcd, spotno, createtime
	)
	SELECT 
		id as adminflightid, carriercd, flightno, crafttype, orgnairportcd, destairportcd, viaairportcd, scheduletime, actualtime, terminal, swing, status, gatecd, checkincounter, exitcd, spotno, '%s' as createtime
	FROM 
		adsb.hnd_flight WHERE id = %v`
	createtime := time.Now().Format("2006-01-02 15:04:05")
	sql = fmt.Sprintf(sql, createtime, adminFlightID)

	lastInsertID, rowsAffected, err = db.Execute(sql)

	return
}
