package utility

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DataHubTagListPlugin() (tagList []interface{}) {
	datahubUrl := G_CONFIGER.GetString("datahub.url")
	Logger.Debug("datahub.url : ", datahubUrl)
	nodeId := G_CONFIGER.GetString("datahub.nodeId")
	Logger.Debug("datahub.nodeid : ", nodeId)
	deviceId := G_CONFIGER.GetString("datahub.deviceId")
	Logger.Debug("datahub.deviceid : ", deviceId)
	apiURL := fmt.Sprintf("%s/v1/Tags/list/"+nodeId+"/"+deviceId, datahubUrl)
	req, _ := http.NewRequest("GET", apiURL, nil)
	datahubInitHeader(req)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		Logger.Error(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		msg := fmt.Sprintf("[%s]: Invalid Status Code. Details [%s]", apiURL, string(body))
		Logger.Error(msg)
	} else {

		var jsonObj map[string]interface{}
		json.Unmarshal([]byte(body), &jsonObj)

		if len(jsonObj) <= 0 {
			msg := fmt.Sprintf("[%s]: Empty Data. Details [%s]", apiURL, string(body))
			Logger.Error(msg)
		} else {
			lists := jsonObj["list"].([]interface{})
			for _, doc := range lists {

				list := doc.(map[string]interface{})

				tagName := list["tagName"].(string)

				tagList = append(tagList, tagName)
			}
		}
	}
	return
}

func datahubInitHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Contect-Type", "application/json")
	ssoToken := G_CONFIGER.GetString("sso.token")
	Logger.Debug("sso.token :", ssoToken)
	req.Header.Set("Authorization", "Bearer "+ssoToken)
}
