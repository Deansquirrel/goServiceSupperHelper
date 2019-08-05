package goServiceSupportHelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper/global"
	"github.com/Deansquirrel/goServiceSupportHelper/object"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type heartBeat struct {
}

func NewHeartBeat() *heartBeat {
	return &heartBeat{}
}

func (ht *heartBeat) HeartBeatUpdate() error {
	if strings.Trim(global.HttpAddress, " ") == "" {
		return errors.New("HttpAddress is empty")
	}
	d := object.HeartBeatRequest{
		ClientId:        global.ClientId,
		HeartBeatClient: time.Now(),
	}
	bd, err := json.Marshal(d)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/heartbeat/update", global.HttpAddress),
		"application/json",
		bytes.NewReader(bd))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var rd object.HeartBeatResponse
	err = json.Unmarshal(body, &rd)
	if err != nil {
		return err
	}
	if rd.ErrCode != 200 {
		return errors.New(rd.ErrMsg)
	}
	if rd.IsForbidden != 0 {
		go func() {
			log.Warn(fmt.Sprintf("Forbidden: %s", rd.ForbiddenReason))
			if global.Cancel != nil {
				global.Cancel()
				time.Sleep(time.Minute)
				os.Exit(0)
			}
		}()
	}
	return nil
}
