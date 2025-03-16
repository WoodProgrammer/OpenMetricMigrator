package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/WoodProgrammer/prom-migrator/cmd"
	prom "github.com/WoodProgrammer/prom-migrator/cmd"
	source "github.com/WoodProgrammer/prom-migrator/lib"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var (
	promHost   string
	promPort   string
	startStamp string
	endStamp   string
	query      string
	step       string
	dataDir    string
	targetDir  string
	metricType string
)

func newPrometheusHandler(host string) prom.Prometheus {
	return &prom.PromHandler{}
}

func CallPrometheus() {
	ch := make(chan string)
	var promHandler prom.Prometheus

	promHandler = newPrometheusHandler(promHost)
	rawMetricData := []string{}

	if !strings.Contains(query, "\"") {
		query = strings.ReplaceAll(query, "{", "{job=\"") // Example fix
	}
	url := fmt.Sprintf("http://%s:%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s",
		promHost, promPort, query, startStamp, endStamp, step)
	status, data := promHandler.FetchPrometheusData(url)

	if status != 200 {
		err := fmt.Errorf("un expected response from Prometheus server %d", status)
		log.Err(err).Msg("Error while running promHandler.FetchPrometheusData")
		return
	}
	rawMetricData = append(rawMetricData, fmt.Sprintf("# TYPE %s %s", strings.Split(query, "{")[0], metricType))

	parsedData, ok := data["data"].(map[string]interface{})

	if !ok {
		err := fmt.Errorf("json parsing error on rawPrometheus data ")
		log.Err(err).Msg("Error parsing 'data'")
		return
	}

	results, ok := parsedData["result"].([]interface{})
	if !ok {
		err := fmt.Errorf("json parsing error on parsedData['result'] data ")
		log.Err(err).Msg("Error parsing 'result'")
		return
	}

	for _, r := range results {
		go promHandler.ParsePrometheusMetric(r, ch)
		tmpData := <-ch
		rawMetricData = append(rawMetricData, tmpData)

	}
	rawMetricData = append(rawMetricData, "# EOF")

	err := ensureDir(dataDir)
	if err != nil {
		log.Err(err).Msg("Error in ensureDir method")
	}

	fileName := fmt.Sprintf("%s/data-%s", dataDir, startStamp)

	if len(targetDir) != 0 {
		cmd.FileHandler(fileName, rawMetricData)
		promHandler.ImportPrometheusData(fileName, targetDir)
	}
}

func main() {

	var rootCmd = &cobra.Command{
		Use:   "openmetricmigrator",
		Short: "CLI tool to export Prometheus data in OpenMetrics format",
		Run: func(cmd *cobra.Command, args []string) {
			CallPrometheus()
		},
	}
	rootCmd.Flags().StringVarP(&promHost, "host", "H", "localhost", "Prometheus host")
	rootCmd.Flags().StringVarP(&promPort, "port", "P", "9090", "Prometheus port")
	rootCmd.Flags().StringVarP(&startStamp, "start", "s", "0", "Start timestamp (epoch)")
	rootCmd.Flags().StringVarP(&endStamp, "end", "e", "0", "End timestamp (epoch)")
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "PromQL query")
	rootCmd.Flags().StringVarP(&step, "step", "t", "15s", "Query step")
	rootCmd.Flags().StringVarP(&dataDir, "directory", "d", "data", "Data directory to export")
	rootCmd.Flags().StringVarP(&targetDir, "targetdir", "T", "", "Target prometheus data directory")
	rootCmd.Flags().StringVarP(&metricType, "metrictype", "m", "counter", "Type of metrics like counter|gauge")

	rootCmd.MarkFlagRequired("start")
	rootCmd.MarkFlagRequired("end")
	rootCmd.MarkFlagRequired("query")

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("CLI execution failed")
		os.Exit(1)
	}
}

func ensureDir(dirName string) error {

	err := os.MkdirAll(dirName, source.DirMode)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
