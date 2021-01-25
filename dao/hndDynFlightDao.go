package dao

import (
	"fmt"
	"hda/db"
	"hda/dto"
)

// SaveHndDynFlight 保存羽田机场动态航班信息
func SaveHndDynFlight(hndDynFlightDto *dto.HndDynFlightDto) (lastInsertID, rowsAffected int64, err error) {

	sql := `
	INSERT INTO adsb.hnd_dynflight_original(
		airlinecd, flightno, 
		orgnairportcd, orgndirectioncd, orgndirectionjpname, orgndirectionenname, 
		destairportcd, destdirectioncd, destdirectionjpname, destdirectionenname, 
		viaairportcd, viadirectioncd, viadirectionjpname, viadirectionenname, 
		scheduletime, actualtime, status, terminal, swing, 
		remarkjpname, remarkenname, remarkjp, remarken, remarkko, remarkhans, remarkhant, 
		fliker, gatecd, remarkcd, checkincounter, spotno, crafttype, operatingstatus, createtime)
	VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %s, %s, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')`

	scheduleTime := hndDynFlightDto.ScheduleTime
	actualTime := hndDynFlightDto.ActualTime

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

	sql = fmt.Sprintf(
		sql,
		hndDynFlightDto.AirlineCD,
		hndDynFlightDto.FlightNo,
		hndDynFlightDto.OrgnAirportCD,
		hndDynFlightDto.OrgnDirectionCD,
		hndDynFlightDto.OrgnDirectionJPName,
		hndDynFlightDto.OrgnDirectionENName,
		hndDynFlightDto.DestAirportCD,
		hndDynFlightDto.DestDirectionCD,
		hndDynFlightDto.DestDirectionJPName,
		hndDynFlightDto.DestDirectionENName,
		hndDynFlightDto.ViaAirportCD,
		hndDynFlightDto.ViaDirectionCD,
		hndDynFlightDto.ViaDirectionJPName,
		hndDynFlightDto.ViaDirectionENName,
		scheduleTime,
		actualTime,
		hndDynFlightDto.Status,
		hndDynFlightDto.Terminal,
		hndDynFlightDto.Swing,
		hndDynFlightDto.RemarkJPName,
		hndDynFlightDto.RemarkENName,
		hndDynFlightDto.RemarkJP,
		hndDynFlightDto.RemarkEN,
		hndDynFlightDto.RemarkKO,
		hndDynFlightDto.RemarkHans,
		hndDynFlightDto.RemarkHant,
		hndDynFlightDto.Fliker,
		hndDynFlightDto.GateCD,
		hndDynFlightDto.RemarkCD,
		hndDynFlightDto.CheckinCounter,
		hndDynFlightDto.SpotNo,
		hndDynFlightDto.CraftType,
		hndDynFlightDto.OperatingStatus,
		hndDynFlightDto.Createtime,
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

	if err != nil {
		return 0, err
	}

	return
}
