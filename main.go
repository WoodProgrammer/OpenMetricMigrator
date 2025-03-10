package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/WoodProgrammer/prom-migrator/cmd"
	prom "github.com/WoodProgrammer/prom-migrator/cmd"
	source "github.com/WoodProgrammer/prom-migrator/lib"

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
)

func CallPrometheus() {

	newPrometheusClient := prom.PromClient{}
	rawMetricData := []string{}

	if !strings.Contains(query, "\"") {
		query = strings.ReplaceAll(query, "{", "{job=\"") // Example fix
	}

	url := fmt.Sprintf("http://%s:%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s",
		promHost, promPort, query, startStamp, endStamp, step)

	data := newPrometheusClient.FetchPrometheusData(url)
	rawMetricData = append(rawMetricData, fmt.Sprintf("# TYPE %s counter", strings.Split(query, "{")[0]))

	for _, result := range data.Data.Result {
		for _, k := range result.Values {
			line := fmt.Sprintf("%s %f %s", query, k[0], k[1])
			rawMetricData = append(rawMetricData, line)
		}
	}
	rawMetricData = append(rawMetricData, fmt.Sprintf("# EOF"))

	err := ensureDir(dataDir)
	if err != nil {
		source.LogErrorWithLine(err, "Error in ensureDir method")
	}

	fileName := fmt.Sprintf("%s/data-%s", dataDir, startStamp)
	cmd.FileHandler(fileName, rawMetricData)

}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "promcli",
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

	rootCmd.MarkFlagRequired("start")
	rootCmd.MarkFlagRequired("end")
	rootCmd.MarkFlagRequired("query")

	if err := rootCmd.Execute(); err != nil {
		source.LogErrorWithLine(err, "CLI execution failed")
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
