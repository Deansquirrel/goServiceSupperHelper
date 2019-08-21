package goServiceSupportHelper

import (
	"context"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper/global"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolCron"
	"github.com/Deansquirrel/goToolEnvironment"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSqlHelper"
	"github.com/Deansquirrel/goToolSVRV3"
	"github.com/Deansquirrel/goToolSVRZ5"
	"github.com/kataras/iris/core/errors"
	"strings"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type Params struct {
	HttpAddress   string
	ClientType    string
	ClientVersion string
	Ctx           context.Context
	Cancel        func()
}

func InitParam(p *Params) {
	global.HttpAddress = strings.Trim(p.HttpAddress, " ")
	global.ClientType = strings.Trim(p.ClientType, " ")
	global.ClientVersion = strings.Trim(p.ClientVersion, " ")
	global.Ctx = p.Ctx
	global.Cancel = p.Cancel
	global.ClientId = getClientId()

	waitForClientId()

	//HeartBeat
	go func() {
		for {
			err := goToolCron.AddFunc(
				"HeartBeat",
				global.HeartBeatCron,
				NewJob().FormatSSJob("HeartBeat", jobHeartBeat),
				panicHandle)
			if err == nil {
				break
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}()
	go refreshClientInfo()
	go refreshHostName()
	go refreshInternetIp()
}

func SetHeartBeatCron(spec string) {
	spec = strings.Trim(spec, " ")
	if spec != "" {
		global.HeartBeatCron = spec
	}
}

func SetRefreshSvrInfoCron(spec string) {
	spec = strings.Trim(spec, " ")
	if spec != "" {
		global.RefreshSvrInfoCron = spec
	}
}

func GetVersion() string {
	return global.Version
}

func GetType() string {
	return global.Type
}

type SvrType string

const (
	SVRNONE SvrType = "svrNone"
	SVRV3   SvrType = "svrV3"
	SVRZ5   SvrType = "svrZ5"
)

func SetOtherInfo(dbConfig *goToolMSSql.MSSqlConfig,
	dbType int,
	svrType SvrType) {
	waitForClientId()
	global.DbConfig = dbConfig
	global.DbType = dbType

	if global.DbConfig == nil {
		return
	}
	go refreshDbId(global.DbConfig, global.DbType)
	switch svrType {
	case SVRV3:
		go func() {
			for {
				err := goToolCron.AddFunc(
					"RefreshSvrV3Info",
					global.RefreshSvrInfoCron,
					NewJob().FormatSSJob("RefreshSvrV3Info", jobRefreshSvrV3Info),
					panicHandle)
				if err == nil {
					break
				} else {
					time.Sleep(time.Second * 10)
				}
			}
		}()
	case SVRZ5:
		go func() {
			for {
				err := goToolCron.AddFunc(
					"RefreshSvrZ5ZlVersion",
					global.RefreshSvrInfoCron,
					NewJob().FormatSSJob("RefreshSvrZ5ZlVersion", jobRefreshSvrZ5ZlVersion),
					panicHandle)
				if err == nil {
					break
				} else {
					time.Sleep(time.Second * 10)
				}
			}
		}()
		go func() {
			for {
				err := goToolCron.AddFunc(
					"RefreshSvrZ5ZlCompany",
					global.RefreshSvrInfoCron,
					NewJob().FormatSSJob("RefreshSvrZ5ZlCompany", jobRefreshSvrZ5ZlCompany),
					panicHandle)
				if err == nil {
					break
				} else {
					time.Sleep(time.Second * 10)
				}
			}
		}()
	default:
		log.Warn(fmt.Sprintf("unexpected type: %s", string(svrType)))
	}
}

func panicHandle(v interface{}) {
	log.Error(fmt.Sprintf("panicHandle: %s", v))
}

func getClientId() string {
	if global.ClientType == "" {
		id := getLocalClientId()
		errMsg := fmt.Sprintf("client type can not be empty,temp Id: %s", id)
		log.Warn(errMsg)
		err := JobErrRecord("GetClientId", errMsg)
		if err != nil {
			log.Error(err.Error())
		}
		return id
	}
	id, err := goToolEnvironment.GetClientId(global.ClientType)
	if err != nil {
		id := getLocalClientId()
		errMsg := fmt.Sprintf("get client id err: %s,temp Id: %s", err.Error(), id)
		log.Warn(errMsg)
		err := JobErrRecord("GetClientId", errMsg)
		if err != nil {
			log.Error(err.Error())
		}
		return id
	}
	return id
}

func getLocalClientId() string {
	u := goToolCommon.Guid()
	u = strings.Replace(u, "-", "", -1)
	return strings.ToUpper(u)
}

//刷新global.InternetIp
func refreshInternetIp() {
	waitForClientId()
	for {
		ip, err := goToolEnvironment.GetInternetAddr()
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		global.InternetIp = ip
		refreshClientInfo()
		break
	}
}

//刷新global.HostName
func refreshHostName() {
	waitForClientId()
	for {
		hostName, err := goToolEnvironment.GetHostName()
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		global.HostName = hostName
		refreshClientInfo()
		break
	}
}

func refreshClientInfo() {
	waitForClientId()
	for {
		dbName := ""
		if global.DbConfig != nil {
			dbName = global.DbConfig.DbName
		}
		err := NewClient().RefreshClientInfo(
			global.ClientId,
			global.ClientType,
			global.ClientVersion,
			global.HostName,
			global.DbId,
			dbName,
			global.InternetIp)
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		break
	}
}

//刷新global.DbId
func refreshDbId(dbConfig *goToolMSSql.MSSqlConfig, dbType int) {
	waitForClientId()
root:
	for {
		switch dbType {
		case 0:
			dbId, err := goToolMSSqlHelper.GetDbId(dbConfig)
			if err != nil {
				time.Sleep(time.Minute * 5)
				continue
			} else {
				global.DbId = dbId
				break root
			}
		case 1:
			dbId, err := goToolMSSqlHelper.GetDbId2000(goToolMSSqlHelper.ConvertDbConfigTo2000(dbConfig))
			if err != nil {
				time.Sleep(time.Minute * 5)
				continue
			} else {
				global.DbId = dbId
				break root
			}
		default:
			return
		}
	}
	refreshClientInfo()
}

func jobHeartBeat(id string) {
	waitForClientId()
	err := NewHeartBeat().HeartBeatUpdate()
	if err != nil {
		log.Error(err.Error())
		err = JobErrRecord(id, err.Error())
		if err != nil {
			log.Error(fmt.Sprintf("record err err: %s", err.Error()))
		}
	}
	return
}

func jobRefreshSvrZ5ZlVersion(id string) {
	waitForClientId()
	vList, err := goToolSVRZ5.GetZlVersion(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		jobHandleErr(id, err)
		return
	}
	if vList == nil {
		jobHandleErr(id, errors.New(fmt.Sprintf("SvrZ5 ZlVersion list is nil")))
		return
	}
	client := NewClient()
	for _, v := range vList {
		err = client.RefreshSvrZ5ZlVersion(global.ClientId,
			v.ObjectName, v.ObjectType, v.ObjectVersion, v.ObjectDate)
		if err != nil {
			jobHandleErr(id, err)
			return
		}
	}
}

func jobRefreshSvrZ5ZlCompany(id string) {
	waitForClientId()
	z, err := goToolSVRZ5.GetZlCompany(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		jobHandleErr(id, err)
		return
	}
	if z == nil {
		jobHandleErr(id, errors.New("zlCompany return nil"))
		return
	}
	err = NewClient().RefreshSvrZ5ZlCompany(global.ClientId,
		z.CoId, z.CoAb, z.CoCode, z.CoType, z.CoUserAb, z.CoUserCode, z.CoAccCrDate)
	if err != nil {
		jobHandleErr(id, err)
		return
	}
}

func jobRefreshSvrV3Info(id string) {
	waitForClientId()
	coId, coAb, coCode, coUserAb, coUserCode, coFunc, err :=
		goToolSVRV3.GetZlCompany(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		jobHandleErr(id, err)
		return
	}
	svName, svVer, svDate, err := goToolSVRV3.GetXtSelfVer(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		jobHandleErr(id, err)
		return
	}
	err = NewClient().RefreshSvrV3Info(
		global.ClientId,
		coId, coAb, coCode, coUserAb, coUserCode, coFunc,
		svName, svVer, svDate)
	if err != nil {
		jobHandleErr(id, err)
		return
	}
}

func jobHandleErr(id string, err error) {
	log.Error(err.Error())
	err = JobErrRecord(id, err.Error())
	if err != nil {
		log.Error(fmt.Sprintf("record err err: %s", err.Error()))
	}
}

func waitForClientId() {
	for {
		if global.ClientId == "" {
			time.Sleep(time.Second)
			continue
		}
		break
	}
}
