package exporter

import (
	"strconv"
	"testing"
)

func TestGetNebulaMetricsJsonNewVersion(t *testing.T)  {
	var ans []StatsMetric
	var err error
	if ans, err = getNebulaMetricsJsonNewVersion("127.0.0.1", 7998); err != nil {
		t.Errorf(err.Error())
	} else {
		for _, an := range ans {
			t.Log(an.Name + ":" + strconv.FormatFloat(an.Value, 'f', -1, 64))
		}
	}
}

func TestGetNebulaRocksDBStatsJson(t *testing.T)  {
	var ans []StatsMetric
	var err error
	if ans, err = getNebulaRocksDBStatsJson("127.0.0.1", 7998); err != nil {
		t.Errorf(err.Error())
	} else {
		for _, an := range ans {
			t.Log(an.Name + ":" + strconv.FormatFloat(an.Value, 'f', -1, 64))
		}
	}
}