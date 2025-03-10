package lib

const (
	DirMode = 0777
)

type PrometheusData struct {
	Status string `json:"status"`
	Data   Data
}

type Data struct {
	ResultType string `json:"resultType"`
	Result     []Result
}
type Result struct {
	Metric   string          `json:"__name__"`
	Instance string          `json:"string"`
	Job      string          `json:"job"`
	Values   [][]interface{} `json:"values"`
}
