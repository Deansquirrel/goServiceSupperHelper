package goServiceSupportHelper

//
//func GetVersion() (string, error) {
//	if strings.Trim(global.HttpAddress, " ") == "" {
//		return "", errors.New("HttpAddress is empty")
//	}
//	resp, err := http.Get(fmt.Sprintf("%s/version", global.HttpAddress))
//	if err != nil {
//		return "", err
//	}
//	defer func() {
//		_ = resp.Body.Close()
//	}()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	var d object.VersionResponse
//	err = json.Unmarshal(body, &d)
//	if err != nil {
//		return "", err
//	}
//	if d.ErrCode != 200 {
//		return "", errors.New(d.ErrMsg)
//	}
//	return d.Version, nil
//}
//
//func GetType() (string, error) {
//	if strings.Trim(global.HttpAddress, " ") == "" {
//		return "", errors.New("HttpAddress is empty")
//	}
//	resp, err := http.Get(fmt.Sprintf("%s/type", global.HttpAddress))
//	if err != nil {
//		return "", err
//	}
//	defer func() {
//		_ = resp.Body.Close()
//	}()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	var d object.TypeResponse
//	err = json.Unmarshal(body, &d)
//	if err != nil {
//		return "", err
//	}
//	if d.ErrCode != 200 {
//		return "", errors.New(d.ErrMsg)
//	}
//	return d.Type, nil
//}
