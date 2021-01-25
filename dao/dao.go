package dao

import (
	"fmt"
	cfgloader "hda/config"
	"hda/db"
)

func init() {

	fmt.Println("Load config file.")

	// 加载配置文件
	cfgloader.LoadConfig("config.toml")
	cfgs, err := cfgloader.GetTable("db")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	driver := cfgs.Get("driver")
	dsn := cfgs.Get("dsn")

	err = db.NewConnection("adsb", driver.(string), dsn.(string))

	if err != nil {
		fmt.Printf("Create adsb connection failed. Details:\n %s", err.Error())
		return
	}
}
