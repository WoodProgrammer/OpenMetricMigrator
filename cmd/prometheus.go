package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	lib "github.com/WoodProgrammer/prom-migrator/lib"
)

type Prometheus interface {
	FetchPrometheusData(url string) lib.PrometheusData
	ImportPrometheusData(file, targetDir string) error
	ExecutePromtoolCommand(args ...string) (string, error)
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

func (promClient *PromClient) ImportPrometheusData(file, targetDir string) error {
	result, err := promClient.ExecutePromtoolCommand(file, targetDir)
	if err != nil {
		lib.LogErrorWithLine(err, "Error on promClient.ExecutePromtoolCommand")
		return err
	}
	fmt.Println("The result is ", result)
	return nil
}

func (promClient *PromClient) ExecutePromtoolCommand(args ...string) (string, error) {
	cmd := exec.Command("promtool tsdb create-blocks-from openmetrics", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
		return stderr.String(), fmt.Errorf("command failed with exit code %d: %s", exitCode, stderr.String())
	}

	return stdout.String(), nil
}
