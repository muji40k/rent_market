package metrics

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var DiskIO = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "system_disk_io",
	Help: "Total bytes written by system to the disk",
}, []string{"disk", "state"})

const DIR string = "/sys/block/"
const STAT string = "/sys/block/%v/stat"
const SECTOR_SIZE float64 = 512

func getSectors(entry fs.DirEntry) (float64, float64) {
	var read float64
	var write float64
	var fields []string

	content, err := os.ReadFile(fmt.Sprintf(STAT, entry.Name()))

	if nil == err {
		fields = strings.Fields(string(content))
		read, err = strconv.ParseFloat(fields[2], 64)
	}

	if nil == err {
		write, err = strconv.ParseFloat(fields[6], 64)
	}

	if nil != err {
		panic(err)
	}

	return read, write
}

func DiskUtilizationWatcher() func() {
	return func() {
		entries, err := os.ReadDir(DIR)

		if nil != err {
			panic(err)
		}

		for _, entry := range entries {
			read, write := getSectors(entry)
			DiskIO.WithLabelValues(entry.Name(), "read").
				Set(read * SECTOR_SIZE)
			DiskIO.WithLabelValues(entry.Name(), "write").
				Set(write * SECTOR_SIZE)
		}
	}
}

