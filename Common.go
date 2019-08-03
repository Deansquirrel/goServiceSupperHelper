package goServiceSupportHelper

import (
	"context"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper/global"
	"github.com/Deansquirrel/goToolCron"
	"github.com/Deansquirrel/goToolEnvironment"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSqlHelper"
	"github.com/Deansquirrel/goToolSVRV3"
	"strings"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type Params struct {
	HttpAddress   string
	ClientType    string
	ClientVersion string
	DbConfig      *goToolMSSql.MSSqlConfig
	//数据库类型，0-非2000,1-2000
	DbType          int
	IsSvrV3         bool
	SvrV3AppType    string
	SvrV3ClientType string
	Ctx             context.Context
	Cancel          func()
}

func InitParam(p *Params) {
	global.HttpAddress = strings.Trim(p.HttpAddress, " ")
	global.ClientType = strings.Trim(p.ClientType, " ")
	global.ClientVersion = strings.Trim(p.ClientVersion, " ")
	global.DbType = p.DbType
	global.IsSvrV3 = p.IsSvrV3
	global.DbConfig = p.DbConfig
	global.Ctx = p.Ctx
	global.Cancel = p.Cancel
	go func() {
		if p.DbConfig == nil {
			return
		}
		refreshDbId(global.DbConfig, global.DbType)
	}()
	go refreshHostName()
	go refreshInternetIp()
	go refreshClientId()
	global.HasInit = true
}

func Start() {
	//检查初始化状态
	for {
		if global.HasInit {
			break
		} else {
			time.Sleep(time.Minute)
			log.Warn(fmt.Sprintf("ServiceSupportHelper 参数设置未完成"))
		}
	}
	if global.HttpAddress == "" {
		log.Warn(fmt.Sprintf("ServiceSupportHelper http address is empty"))
		return
	}
	for {
		if global.ClientId == "" {
			time.Sleep(time.Minute)
			continue
		} else {
			break
		}
	}
	go func() {
		for {
			err := goToolCron.AddFunc(
				"HeartBeatUpdate",
				"0 * * * * ?",
				FormatSSJob("HeartBeatUpdate", jobHeartBeatUpdate),
				panicHandle)
			if err != nil {
				log.Error(err.Error())
				time.Sleep(time.Minute)
				continue
			} else {
				break
			}
		}
	}()
	go func() {
		for {
			if global.ClientId == "" {
				time.Sleep(time.Minute)
				continue
			}
			ip := global.InternetIp
			err := RefreshFlashInfo(global.ClientId, global.Version, ip)
			if err != nil {
				log.Error(err.Error())
				time.Sleep(time.Minute * 10)
				continue
			}
			if ip != "" {
				break
			} else {
				time.Sleep(time.Minute)
			}
		}
	}()
	go func() {
		if global.DbConfig != nil && global.IsSvrV3 {
			for {
				err := goToolCron.AddFunc(
					"RefreshSvrV3Info",
					"0 15/30 * * * ?",
					FormatSSJob("RefreshSvrV3Info", jobRefreshSvrV3Info),
					panicHandle)
				if err != nil {
					log.Error(err.Error())
					time.Sleep(time.Minute)
					continue
				} else {
					break
				}
			}
		}
	}()
}

//Demo
//func init(){
//	//goServiceSupportHelper.HttpAddress = "http://192.168.8.148:8000"
//	goServiceSupportHelper.InitParam(&goServiceSupportHelper.Params{
//		HttpAddress:"http://192.168.8.148:8000",
//		ClientType:global.Type,
//		ClientVersion:global.Version,
//		DbConfig:&goToolMSSql.MSSqlConfig{
//			Server:"192.168.5.1",
//			Port:2003,
//			User:"sa",
//			Pwd:"",
//			DbName:"Z9门店",
//		},
//		//数据库类型，0-非2000,1-2000
//		DbType:1,
//		IsSvrV3:true,
//		SvrV3AppType:"83",
//		SvrV3ClientType:"8301",
//	})
//	go goServiceSupportHelper.Start()
//}

//刷新global.ClientId
func refreshClientId() {
	dbId := -1
	dbName := ""
	for {
		if global.HostName == "" {
			time.Sleep(time.Minute)
			continue
		}
		if global.DbConfig == nil {
			break
		}
		if global.DbId < 1 {
			time.Sleep(time.Minute)
			continue
		}
		dbId = global.DbId
		dbName = global.DbConfig.DbName
		break
	}
	for {
		clientId, err := GetClientId(global.ClientType, global.HostName, dbId, dbName)
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		global.ClientId = clientId
		break
	}
	return
}

//刷新global.InternetIp
func refreshInternetIp() {
	for {
		ip, err := goToolEnvironment.GetInternetAddr()
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		global.InternetIp = ip
		return
	}
}

//刷新global.HostName
func refreshHostName() {
	for {
		hostName, err := goToolEnvironment.GetHostName()
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		global.HostName = hostName
		return
	}
}

//刷新global.DbId
func refreshDbId(dbConfig *goToolMSSql.MSSqlConfig, dbType int) {
	for {
		switch dbType {
		case 0:
			dbId, err := goToolMSSqlHelper.GetDbId(dbConfig)
			if err != nil {
				time.Sleep(time.Minute * 5)
				continue
			} else {
				global.DbId = dbId
				return
			}
		case 1:
			dbId, err := goToolMSSqlHelper.GetDbId2000(goToolMSSqlHelper.ConvertDbConfigTo2000(dbConfig))
			if err != nil {
				time.Sleep(time.Minute * 5)
				continue
			} else {
				global.DbId = dbId
				return
			}
		default:
			return
		}
	}
}

func panicHandle(v interface{}) {
	log.Error(fmt.Sprintf("panicHandle: %s", v))
}

func jobHeartBeatUpdate() {
	err := HeartBeatUpdate(global.ClientId)
	if err != nil {
		log.Error(err.Error())
	}
}

func jobRefreshSvrV3Info() {
	coId, coAb, coCode, coUserAb, coUserCode, coFunc, err :=
		goToolSVRV3.GetZlCompany(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		log.Error(err.Error())
		return
	}
	svName, svVer, svDate, err := goToolSVRV3.GetXtSelfVer(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = RefreshSvrV3Info(
		global.ClientId,
		coId, coAb, coCode, coUserAb, coUserCode, coFunc,
		svName, svVer, svDate)
	if err != nil {
		log.Error(err.Error())
		return
	}
}
