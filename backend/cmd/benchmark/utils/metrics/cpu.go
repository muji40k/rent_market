package metrics

import (
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var CpuUtilization = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "system_cpu_utilization",
	Help: "Cpu utilization in %",
})

func getStats() (uint64, uint64) {
	contents, err := os.ReadFile("/proc/stat")

	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)

		if fields[0] != "cpu" {
			continue
		}

		var idle uint64
		var other uint64
		numFields := len(fields)

		for i := 1; i < numFields; i++ {
			val, err := strconv.ParseUint(fields[i], 10, 64)
			if err != nil {
				panic(err)
			}

			if i == 4 {
				idle = val
			} else {
				other += val
			}
		}

		return other, idle
	}

	panic("No cpu entry found")
}

func CPUUtilizationWatcher() func() {
	pother, pidle := getStats()

	return func() {
		other, idle := getStats()

		didle := idle - pidle
		dother := other - pother
		pidle = idle
		pother = other

		CpuUtilization.Set(float64(dother) / float64(dother+didle) * 100)
	}
}

