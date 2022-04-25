package exporter

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	extraLabelNames []string
	extraLabelValues []string
)

type StatsMetric struct {
	Name string 	`json:"name"`
	Value float64		`json:"value"`
}

func InitExtraLabels(config StaticConfig) {
	if config.ExtraLabels != nil {
		for _, extraLabel := range config.ExtraLabels {
			extraLabelNames = append(extraLabelNames, extraLabel.Name)
			extraLabelValues = append(extraLabelValues, extraLabel.Value)
		}
	}
}

func getNebulaMetrics(ipAddress string, port int32) ([]string, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/stats", ipAddress, port))
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	metrics := strings.Split(strings.TrimSpace(string(bytes)), "\n")

	return metrics, nil
}

func getNebulaRocksDBStats(ipAddress string, port int32) ([]string, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/rocksdb_stats", ipAddress, port))
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	metrics := strings.Split(strings.TrimSpace(string(bytes)), "\n")

	return metrics, nil
}

func isNebulaComponentRunning(ipAddress string, port int32) bool {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/status?format=json", ipAddress, port))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	type nebulaStatus struct {
		GitInfoSha string `json:"git_info_sha"`
		Status     string `json:"status"`
	}

	var status nebulaStatus
	if err := json.Unmarshal(bytes, &status); err != nil {
		return false
	}

	return status.Status == "running"
}

func getNebulaMetricsJson(ipAddress string, port int32) ([]StatsMetric, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/stats?format=json", ipAddress, port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var metrics []StatsMetric
	if err := json.Unmarshal(bytes, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func getNebulaRocksDBStatsJson(ipAddress string, port int32) ([]StatsMetric, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/rocksdb_stats?format=json", ipAddress, port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var metrics []map[string]string
	if err := json.Unmarshal(bytes, &metrics); err != nil {
		return nil, err
	}

	var statsMetric []StatsMetric
	for _, metric := range metrics {
		v, err := strconv.ParseFloat(metric["value"], 64)
		if err != nil {
			continue
		}
		statsMetric = append(statsMetric, StatsMetric{
			metric["name"],
			v,
		})
	}

	return statsMetric, nil
}


func mustNewConstMetric(desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) prometheus.Metric {
	if labelValues != nil {
		labelValues = append(labelValues, extraLabelValues...)
	} else {
		labelValues = extraLabelValues
	}

	metric := prometheus.MustNewConstMetric(desc, valueType, value, labelValues...)

	return metric
}

func newDesc(fqName string, docString string, labelNames ...string) *prometheus.Desc {
	if labelNames != nil {
		labelNames = append(labelNames, extraLabelNames...)
	} else {
		labelNames = extraLabelNames
	}

	return prometheus.NewDesc(
		fqName,
		docString,
		labelNames,
		nil)
}