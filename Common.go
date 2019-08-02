package goServiceSupportHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper/object"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	HttpAddress = ""
)

//==============================================================================

//base

func GetVersion() (string, error) {
	if strings.Trim(HttpAddress, " ") == "" {
		return "", errors.New("HttpAddress is empty")
	}
	resp, err := http.Get(fmt.Sprintf("%s/version", HttpAddress))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var d object.VersionResponse
	err = json.Unmarshal(body, &d)
	if err != nil {
		return "", err
	}
	if d.ErrCode != 200 {
		return "", errors.New(d.ErrMsg)
	}
	return d.Version, nil
}

func GetType() (string, error) {
	if strings.Trim(HttpAddress, " ") == "" {
		return "", errors.New("HttpAddress is empty")
	}
	resp, err := http.Get(fmt.Sprintf("%s/type", HttpAddress))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var d object.TypeResponse
	err = json.Unmarshal(body, &d)
	if err != nil {
		return "", err
	}
	if d.ErrCode != 200 {
		return "", errors.New(d.ErrMsg)
	}
	return d.Type, nil
}

//==============================================================================

//client

func GetClientId() (string, error) {
	//TODO
	return "", nil
}

func RefreshFlashInfo() error {
	//TODO
	return nil
}

func RefreshSvrV3Info() error {
	//TODO
	return nil
}

//==============================================================================

//heartbeat

func HeartBeatUpdate() error {
	//TODO
	return nil
}

//==============================================================================

//job

func JobRecordStart() error {
	//TODO
	return nil
}

func JobRecordEnd() error {
	//TODO
	return nil
}

//==============================================================================
