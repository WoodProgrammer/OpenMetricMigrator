package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type Prometheus interface {
	FetchPrometheusData(url string) (int, map[string]interface{})
	ImportPrometheusData(file, targetDir string) error
	ParsePrometheusMetric(r interface{}, ch chan interface{}) []string
	ExecutePromtoolCommand(sourceDir, targetDir string) (string, error)
}

type PromHandler struct {
}

func (promHandler *PromHandler) FetchPrometheusData(url string) (int, map[string]interface{}) {
	var metric map[string]interface{}
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Err(err).Msg("Failed to fetch data")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("Error while reading response")
		return resp.StatusCode, metric
	}

	err = json.Unmarshal([]byte(body), &metric)
	if err != nil {
		log.Err(err).Msg("Error while marshaling metric data")
		return resp.StatusCode, metric
	}

	return resp.StatusCode, metric
}

func (promHandler *PromHandler) ImportPrometheusData(file, targetDir string) error {
	result, err := promHandler.ExecutePromtoolCommand(file, targetDir)
	if err != nil {
		log.Err(err).Msg("Error on PromHandler.ExecutePromtoolCommand")
		return err
	}
	log.Info().Msgf("Promtool command output is %s", result)
	return nil
}

func (promHandler *PromHandler) ExecutePromtoolCommand(sourceDir, targetDir string) (string, error) {
	cmd := exec.Command("promtool", "tsdb", "create-blocks-from", "openmetrics", sourceDir, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Err(err).Msg("cmd.Run() failed with \n")
	}
	return string(output), err
}

func (promHandler *PromHandler) ParsePrometheusMetric(r interface{}, ch chan interface{}) []string {
	rawMetricData := []string{}
	result, _ := r.(map[string]interface{})

	labelMap := []string{}
	metric, ok := result["metric"].(map[string]interface{})
	metricName, ok := metric["__name__"].(string)
	if ok {
		for key, value := range metric {

			if key != "__name__" {
				labelMap = append(labelMap, fmt.Sprintf(`%s="%s"`, key, value))
			}
		}
		query := fmt.Sprintf(`%s{%s}`, metricName, strings.Join(labelMap, ","))

		values, ok := result["values"].([]interface{})
		if ok {
			for _, v := range values {
				valArr, ok := v.([]interface{})
				if ok && len(valArr) == 2 {
					tmpData := fmt.Sprintf("%s %v %f", query, valArr[1], valArr[0])
					rawMetricData = append(rawMetricData, tmpData)
				}
			}
		}
	}
	return rawMetricData
}
