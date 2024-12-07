package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var NsPerOp = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Name: "benchmark_ns_per_op",
	Help: "BENCHMARK: Nanoseconds took to proceed operation",
	Objectives: map[float64]float64{
		0.5:  0,
		0.9:  0,
		0.95: 0.,
	},
}, []string{"target", "target_task"})

var AllocsPerOp = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Name: "benchmark_allocs_per_op",
	Help: "BENCHMARK: Amount of allocations during operation",
	Objectives: map[float64]float64{
		0.5:  0,
		0.9:  0,
		0.95: 0.,
	},
}, []string{"target", "target_task"})

var AllocedBytesPerOp = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Name: "benchmark_alloced_bytes_per_op",
	Help: "BENCHMARK: Amount of allocated bytes during operation",
	Objectives: map[float64]float64{
		0.5:  0,
		0.9:  0,
		0.95: 0.,
	},
}, []string{"target", "target_task"})

const CALL_PER_OP string = "callns/op"

var CallDuration = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "benchmark_call_time",
	Help: "BENCHMARK: Nanoseconds for target call",
}, []string{"target", "target_task"})

const SERIALIZATION_PER_OP string = "serns/op"

var SerializationDuration = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "benchmark_serialization_time",
	Help: "BENCHMARK: Nanoseconds took to serialize collection",
}, []string{"target", "target_task"})

