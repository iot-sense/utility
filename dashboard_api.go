package utility

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

//DashboardUpdateDashboard func
func DashboardUpdateDashboard(content []byte) {
	apiHost := G_CONFIGER.GetString("dashboard.host")
	apiURL := fmt.Sprintf("%s/api/dashboards/db", apiHost)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(content))
	dashboardInitHeader(req)

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
		Logger.Info(string(body))
	}
}

//DashboardGetDashboardByUid func
func DashboardGetDashboardByUid(uid string) (content []byte) {
	apiHost := G_CONFIGER.GetString("dashboard.host")
	apiURL := fmt.Sprintf("%s/api/dashboards/uid/%s", apiHost, uid)
	req, _ := http.NewRequest("GET", apiURL, nil)
	dashboardInitHeader(req)

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
		content = body
	}

	return content
}

//dashboardInitHeader func
func dashboardInitHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	ssoToken := G_CONFIGER.GetString("sso.token")
	req.Header.Set("Authorization", "Bearer "+ssoToken)
}
