package router

import (
	"bytes"
	"fmt"
	"net/http"
	"rent_service/cmd/benchmark/utils"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const ITER string = "iter"

func BenchmarkRunnerRoute(bench utils.Benchmark, reporter utils.Reporter) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		iter, err := strconv.ParseUint(ctx.Params.ByName(ITER), 10, 64)

		if nil != err {
			ctx.String(http.StatusInternalServerError, err.Error())
		} else {
			var writer = bytes.Buffer{}
			start := time.Now()
			writer.WriteString(fmt.Sprintln(start))

			utils.RunBenchmark(bench, uint(iter),
				utils.ReporterFunc(func(res *testing.BenchmarkResult) {
					if nil != reporter {
						reporter.Report(res)
					}

					res.Extra[".B/op"] = float64(res.AllocedBytesPerOp())
					res.Extra[".allocs/op"] = float64(res.AllocsPerOp())

					fmt.Fprintf(&writer, "%v: %v\n", time.Now(), res.String())
				}),
			)

			end := time.Now()
			writer.WriteString(fmt.Sprintln(end))
			writer.WriteString(fmt.Sprintln(end.Sub(start)))

			ctx.Data(http.StatusOK, "text", writer.Bytes())
		}
	}
}

