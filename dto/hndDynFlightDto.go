package dto

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// HndFlightDto 羽田航班信息
type HndFlightDto struct {
	ID             int64
	CarrierCD      string
	FlightNo       string
	CraftType      string
	OrgnAirportCd  string
	DestAirportCd  string
	ViaAirportCd   string
	ScheduleTime   string
	ActualTime     string
	Terminal       string
	Swing          string
	Status         string
	GateCd         string
	CheckinCounter string
	ExitCD         string
	SpotNo         string
	CreateTime     string
}

// HashValue 计算航班信息对象的hash值
func (f *HndFlightDto) HashValue() string {

	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%v", *f))
	return fmt.Sprint(h.Sum(nil))
}

// HndShareCodeDto sharecode flight info
type HndShareCodeDto struct {
	ID            int
	AdminFlightID int64
	AirlineCD     string
	FlightNo      string
}
