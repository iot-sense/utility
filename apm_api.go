package utility

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//APMPlugin struct
type APMPlugin struct {
	Category string `json:"category"`
	URL      string `json:"url"`
}

//APMRegPlugin func
func APMRegPlugin(apiList []APMPlugin) {
	apmHost := G_CONFIGER.GetString("apm.host")
	Logger.Debug("apm.host :", apmHost)
	apiURL := fmt.Sprintf("%s/measure/plugin/register", apmHost)
	payload, _ := json.Marshal(apiList)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	apmInitHeader(req)

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
	}
}

//APMUnregPlugin func
func APMUnregPlugin(categoryList []string) {
	apmHost := G_CONFIGER.GetString("apm.host")
	apiURL := fmt.Sprintf("%s/measure/plugin/unregister", apmHost)
	payload, _ := json.Marshal(categoryList)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	apmInitHeader(req)

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
	}
}

//APMListPlugin func
func APMListPlugin() (pluginList []APMPlugin) {
	apmHost := G_CONFIGER.GetString("apm.host")
	apiURL := fmt.Sprintf("%s/measure/plugin/url", apmHost)
	req, _ := http.NewRequest("GET", apiURL, nil)
	apmInitHeader(req)

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
		json.Unmarshal([]byte(body), &pluginList)
		if len(pluginList) <= 0 {
			msg := fmt.Sprintf("[%s]: Empty Data. Details [%s]", apiURL, string(body))
			Logger.Error(msg)
		}
	}

	return
}

//apmInitHeader func
func apmInitHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	ssoToken := G_CONFIGER.GetString("sso.token")
	req.Header.Set("Authorization", "Bearer "+ssoToken)
	apmTenantID := G_CONFIGER.GetString("apm.tenant_id")
	req.Header.Set("spaceid", apmTenantID)
}
