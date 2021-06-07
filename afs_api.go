package utility

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//AFSGetInferenceUrl func
func AFSGetInferenceUrl() (inferenceMap map[string]string) {
	apiHost := G_CONFIGER.GetString("afs.host")
	instanceId := G_CONFIGER.GetString("afs.instanceId")
	apiURL := fmt.Sprintf("%s/v2/instances/%s/cloud_inferences", apiHost, instanceId)
	req, _ := http.NewRequest("GET", apiURL, nil)
	afsInitHeader(req)

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
		if resources, ok := result["resources"].([]interface{}); ok && len(resources) > 0 {
			inferenceMap = map[string]string{}
			for _, resource := range resources {
				r := resource.(map[string]interface{})
				name := r["name"].(string)
				url := r["url"].(string)
				inferenceMap[name] = url
			}
		}
	}
	Logger.Info(inferenceMap)
	return
}

//afsInitHeader func
func afsInitHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	ssoToken := G_CONFIGER.GetString("sso.token")
	req.Header.Set("Authorization", "Bearer "+ssoToken)
}
