package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Prometheus interface {
	FetchPrometheusData(url string) (int, map[string]interface{})
	ImportPrometheusData(file, targetDir string) error
	ExecutePromtoolCommand(args ...string) (string, error)
}

type PromClient struct {
}

func (promClient *PromClient) FetchPrometheusData(url string) (int, map[string]interface{}) {
	var metric map[string]interface{}

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

func (promClient *PromClient) ImportPrometheusData(file, targetDir string) error {
	result, err := promClient.ExecutePromtoolCommand(file, targetDir)
	if err != nil {
		log.Err(err).Msg("Error on promClient.ExecutePromtoolCommand")
		return err
	}
	log.Info().Msgf("Promtool command output is %s", result)
	return nil
}

func (promClient *PromClient) ExecutePromtoolCommand(sourceDir, targetDir string) (string, error) {
	cmd := exec.Command("promtool", "tsdb", "create-blocks-from", "openmetrics", sourceDir, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Err(err).Msg("cmd.Run() failed with \n")
	}
	return string(output), err
}
