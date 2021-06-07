package utility

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//SSOGetAuth func
func SSOGetAuth(uid string, pwd string, token string) (ssoToken string) {
	ssoToken = ssoLoginByToken(token)
	if ssoToken == "" {
		ssoToken = ssoLoginByPwd(uid, pwd)
	}
	return
}

//ssoInitHeader func
func ssoInitHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}

//ssoLoginByToken func
func ssoLoginByToken(token string) (ssoToken string) {
	ssoToken = ""
	apiHost := G_CONFIGER.GetString("sso.host")
	apiURL := fmt.Sprintf("%s/tokenvalidation", apiHost)
	payload, _ := json.Marshal(map[string]string{
		"token": token,
	})
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	ssoInitHeader(req)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		Logger.Error(err)
		return ssoToken
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		msg := fmt.Sprintf("[%s]: Invalid Status Code %d. Details [%s]", apiURL, res.StatusCode, string(body))
		Logger.Error(msg)
		return ssoToken
	} else {
		var result map[string]interface{}
		json.Unmarshal([]byte(body), &result)
		if accessToken, ok := result["accessToken"].(string); ok {
			ssoToken = accessToken
		}
	}
	return ssoToken
}

//ssoLoginByPwd func
func ssoLoginByPwd(uid string, pwd string) (ssoToken string) {
	apiHost := G_CONFIGER.GetString("sso.host")
	apiURL := fmt.Sprintf("%s/auth/native", apiHost)
	payload, _ := json.Marshal(map[string]string{
		"username": uid,
		"password": pwd,
	})
	Logger.Debug("username", uid, "password", pwd)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	ssoInitHeader(req)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		Logger.Error(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		msg := fmt.Sprintf("[%s]: Invalid Status Code %d. Details [%s]", apiURL, res.StatusCode, string(body))
		Logger.Error(msg)
	} else {
		var result map[string]interface{}
		json.Unmarshal([]byte(body), &result)
		if accessToken, ok := result["accessToken"].(string); ok {
			ssoToken = accessToken
		}
	}
	return
}
