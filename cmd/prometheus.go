package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	lib "github.com/WoodProgrammer/prom-migrator/lib"
)

type Prometheus interface {
	FetchPrometheusData(url string) lib.PrometheusData
}

type PromClient struct {
}

func (promClient *PromClient) FetchPrometheusData(url string) lib.PrometheusData {
	var metric lib.PrometheusData

	resp, err := http.Get(url)
	if err != nil {
		lib.LogErrorWithLine(err, "Failed to fetch data")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		lib.LogErrorWithLine(err, "Error while reading response")
		return metric
	}

	err = json.Unmarshal([]byte(body), &metric)
	if err != nil {
		lib.LogErrorWithLine(err, "Error while marshaling metric data")
		return metric
	}

	return metric
}
