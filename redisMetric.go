package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type redisMetric struct {
	name  string
	value float64
}

func generateUniqueMetrics(reIn string) []redisMetric {
	var metrics []redisMetric

	// ensure that the map is always iterated in the same order
	metrNams := make([]string, 0, len(uniqueMetricMap))
	for k := range uniqueMetricMap {
		metrNams = append(metrNams, k)
	}
	sort.Strings(metrNams)

	for _, metrNam := range metrNams {
		submNam := uniqueMetricMap[metrNam]
		val, err := fetchMetricValue(reIn, metrNam)

		if err == nil {
			metrics = append(metrics, redisMetric{
				name:  submNam,
				value: val,
			})
		} else {
			log.Println(err)
		}

	}
	return metrics
}

func fetchMetricValue(reIn string, metrNam string) (float64, error) {
	scanner := bufio.NewScanner(strings.NewReader(reIn))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, metrNam) {
			value, err := strconv.ParseFloat(strings.Split(line, ":")[1], 64)
			errCheckFatal(err)

			return value, nil
		}
	}
	return 0, fmt.Errorf("metric %s not found", metrNam)
}

func generateRecordsMetrics(reIn string) []redisMetric {
	var metrics []redisMetric
	// group 1 is the db ID; group 2 is the number of keys (records)
	re := regexp.MustCompile(`^db(\d+):keys=(\d+)`)

	scanner := bufio.NewScanner(strings.NewReader(reIn))
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)

		if len(match) > 0 {
			id := match[1]
			value, _ := strconv.ParseFloat(match[2], 64)

			metrics = append(metrics, redisMetric{
				name:  fmt.Sprintf("records-%s", id),
				value: value,
			})
		}

	}
	return metrics
}

func parsePutvalString(instance string, metric redisMetric) string {
	var submitValue string
	if strings.Contains(metric.name, "ps_cputime") {
		//redis returns decimal digits for CPU usage but collectd expects none
		submitValue = strconv.FormatFloat(metric.value, 'f', 0, 64)
	} else {
		// only return the required decimal digitas and trim trailing 0s
		submitValue = strconv.FormatFloat(metric.value, 'f', -1, 64)
	}
	return fmt.Sprintf("PUTVAL \"%s/redis-%s/%s\" interval=%f N:%s\n", hostname, instance, metric.name, collectdInterval, submitValue)
}
