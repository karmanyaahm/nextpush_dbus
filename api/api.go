package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"zgo.at/errors"

	"unifiedpush.org/go/nextpush_dbus/auth"
)

const (
	UserAgent       = "UnifiedPush_NextPush_DBus/1.0"
	NCAppPathPrefix = "/index.php/apps/uppush"
	pathDevice      = NCAppPathPrefix + "/device/"
	pathApp         = NCAppPathPrefix + "/app/"
)

func request(method, url string, body io.Reader) (req *http.Request, err error) {
	server, uname, passwd, err := auth.GetCreds()

	req, err = http.NewRequest(method, server+url, body)
	req.Header.Set("User-Agent", UserAgent)
	req.SetBasicAuth(uname, passwd)
	return
}

func requestDo(method, url string) (resp *http.Response, err error) {
	req, err := request(method, url, nil)
	if err != nil {
		return
	}
	resp, err = http.DefaultClient.Do(req)
	return
}

func requestDoFancy(method, url string, respBody interface{}) error {
	resp, err := requestDo(method, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("SERVER ERROR " + resp.Status + " " + url)
	}

	return readUnmarshal(&respBody, resp.Body)
}

func CreateDevice(name string) (deviceID string, err error) {
	q := url.Values{"deviceName": {name}}
	parseURL := url.URL{Path: pathDevice, RawQuery: q.Encode()}

	respVal := struct {
		Success  bool
		DeviceId string
	}{}
	err = requestDoFancy("PUT", parseURL.String(), &respVal)
	if err != nil {
		return
	} else if !respVal.Success {
		err = errors.New("SERVER ERROR: UNKNOWN")
	}

	deviceID = respVal.DeviceId
	return
}

//The token this returns is NOT the full endpoint
func CreateApp(deviceId, name string) (token string, err error) {
	q := url.Values{"deviceId": {deviceId},
		"appName": {name}}
	parseURL := url.URL{Path: pathApp, RawQuery: q.Encode()}

	respVal := struct {
		Success bool
		Token   string
	}{}
	err = requestDoFancy("PUT", parseURL.String(), &respVal)
	if err != nil {
		return
	} else if !respVal.Success {
		err = errors.New("SERVER ERROR: UNKNOWN")
	}

	token = respVal.Token
	return
}

// caller is responsible for checking auth.GetCreds() at the start of the session to know about any possible errors beforehand
func GetEndpointFromApp(token string) (endpoint string) {
	server, _, _, _ := auth.GetCreds()
	return server + NCAppPathPrefix + "/push/" + token
}

func DeleteDevice(deviceID string) (err error) {
	respVal := struct {
		Success bool
	}{}
	err = requestDoFancy("DELETE", pathDevice+deviceID, &respVal)
	if err != nil {
		return
	} else if !respVal.Success {
		err = errors.New("SERVER ERROR: UNKNOWN")
	}

	return
}

func DeleteApp(token string) (err error) {
	respVal := struct {
		Success bool
	}{}
	err = requestDoFancy("DELETE", pathApp+token, &respVal)
	if err != nil {
		return
	} else if !respVal.Success {
		err = errors.New("SERVER ERROR: UNKNOWN")
	}

	return
}

func readUnmarshal(i interface{}, body io.Reader) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "IO ERROR")
	}
	err = json.Unmarshal(b, i)
	return errors.Wrap(err, "JSON ERROR")
}
