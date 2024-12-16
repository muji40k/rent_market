package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rent_service/cmd/benchmark/utils"
	"rent_service/cmd/benchmark/utils/db"
	"rent_service/cmd/benchmark/utils/metrics"
	product_bench "rent_service/cmd/benchmark/utils/product"
	"rent_service/cmd/benchmark/utils/router"
	gormrepo "rent_service/internal/repository/implementation/gorm/repositories/product"
	sqlxrepo "rent_service/internal/repository/implementation/sql/repositories/product"
	"rent_service/internal/repository/interfaces/product"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getPgxConnectionString(config db.Config) string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}

func getGormConnectionString(config db.Config) string {
	return fmt.Sprintf(
		"host=%v user=%v dbname=%v password=%v port=%v",
		config.Host, config.User, config.Database, config.Password,
		config.Port,
	)
}

func getReporter(realisation string, method string) utils.Reporter {
	return utils.ReporterFunc(func(res *testing.BenchmarkResult) {
		if callTime, found := res.Extra[metrics.CALL_PER_OP]; found {
			metrics.CallDuration.WithLabelValues(realisation, method).
				Set(callTime)
		}

		if serTime, found := res.Extra[metrics.SERIALIZATION_PER_OP]; found {
			metrics.SerializationDuration.WithLabelValues(realisation, method).
				Set(serTime)
		}

		metrics.NsPerOp.WithLabelValues(realisation, method).
			Observe(float64(res.NsPerOp()))
		metrics.AllocsPerOp.WithLabelValues(realisation, method).
			Observe(float64(res.AllocsPerOp()))
		metrics.AllocedBytesPerOp.WithLabelValues(realisation, method).
			Observe(float64(res.AllocedBytesPerOp()))
	})
}

func initRoutes(engine *gin.Engine, config db.Config) {
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.GET("/call/gorm", func(ctx *gin.Context) {
		db, err := gorm.Open(
			postgres.Open(getGormConnectionString(config)),
			&gorm.Config{},
		)

		if nil != err {
			ctx.Status(http.StatusInternalServerError)
			panic(err)
		}

		defer func() {
			sqldb, _ := db.DB()

			if nil != sqldb {
				sqldb.Close()
			}
		}()

		c, s := product_bench.SingleCall(gormrepo.New(db))
		ctx.Data(http.StatusOK, "text", []byte(fmt.Sprintf("%v %v", c, s)))
	})
	engine.GET("/call/sqlx", func(ctx *gin.Context) {
		db, err := sqlx.Connect("pgx", getPgxConnectionString(config))

		if nil != err {
			panic(err)
		}

		defer func() {
			if nil != db {
				db.Close()
			}
		}()

		c, s := product_bench.SingleCall(sqlxrepo.New(db))
		ctx.Data(http.StatusOK, "text", []byte(fmt.Sprintf("%v %v", c, s)))
	})
	engine.GET(
		fmt.Sprintf("/bench/sqlx/:%v", router.ITER),
		router.BenchmarkRunnerRoute(
			product_bench.New(func() (product.IRepository, func()) {
				db, err := sqlx.Connect("pgx", getPgxConnectionString(config))

				if nil != err {
					panic(err)
				}

				return sqlxrepo.New(db), func() { db.Close() }
			}),
			getReporter("PSQLProductRepository", "GetWithFilter"),
		),
	)
	engine.GET(
		fmt.Sprintf("/bench/gorm/:%v", router.ITER),
		router.BenchmarkRunnerRoute(
			product_bench.New(func() (product.IRepository, func()) {
				db, err := gorm.Open(
					postgres.Open(getGormConnectionString(config)),
					&gorm.Config{},
				)

				if nil != err {
					panic(err)
				}

				return gormrepo.New(db), func() {
					sqldb, _ := db.DB()

					if nil != sqldb {
						sqldb.Close()
					}
				}
			}),
			getReporter("GORMProductRepository", "GetWithFilter"),
		),
	)
}

func main() {
	testing.Init()

	config := db.FromEnv()
	engine := gin.Default()

	initRoutes(engine, config)

	serv := &http.Server{
		Addr:    "0.0.0.0:2112",
		Handler: engine,
	}

	var wg sync.WaitGroup
	var exit = false

	wg.Add(1)
	go func() {
		defer wg.Done()
		watcher := metrics.CPUUtilizationWatcher()

		for !exit {
			time.Sleep(3 * time.Second)
			watcher()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		watcher := metrics.DiskUtilizationWatcher()

		for !exit {
			time.Sleep(3 * time.Second)
			watcher()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	exit = true
	serv.Close()
	wg.Wait()
}

