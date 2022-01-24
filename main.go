package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gammazero/workerpool"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mailru/go-clickhouse"
	"github.com/sirupsen/logrus"
)

var asudb *sql.DB
var chdb *sql.DB
var log = logrus.New()

func getDSN(fn string) (string, string) {
	err := godotenv.Load(os.ExpandEnv(fn))
	if err != nil {
		panic(err.Error())
	}
	tsdb := os.Getenv("TS_DB")
	asudb := os.Getenv("DB_DSN")

	return tsdb, asudb
}

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <YYYYMM> db.conf <dual>(optional)\n", os.Args[0])
		return
	}

	date := os.Args[1]

	chDb, asuDb := getDSN(os.Args[2])
	fmt.Println(chDb, asuDb)

	var dual string
	if len(os.Args) == 4 {
		dual = os.Args[3]
	}

	initLog(date)
	initDB(chDb, asuDb)
	defer chdb.Close()
	defer asudb.Close()
	loadPrefixes()

	var requests []Job
	if dual != "" {
		requests = getJobsDual(date)
	} else {
		requests = getJobs(date)
	}

	if len(requests) == 0 {
		fmt.Println("Unable to get jobs, exiting")
		return
	}
	reestr := &Reestr{}
	reestr.createReestrFiles(date)

	wp := workerpool.New(200)
	for _, r := range requests {
		r := r
		wp.Submit(func() {
			//fmt.Println("main",r.orgId)
			runWork(r.orgId, r.date, r.dual, reestr)
		})
	}

	wp.StopWait()
	reestr.closeFiles()
}

func initDB(chDb, asuDb string) {
	var err error
	//Init clickhouse
	//chdb, err = sql.Open("clickhouse", "http://default:te410pte412p@10.19.64.124:8123/asubill?debug=true")
	chdb, err = sql.Open("clickhouse", chDb)
	if err != nil {
		panic(err.Error())
	}
	//Init ASUDB
	//asudb, err = sql.Open("mysql", "asu:qwerty@tcp(localhost:3308)/agts_asu")
	asudb, err = sql.Open("mysql", asuDb)
	if err != nil {
		panic(err.Error())
	}
}
func initLog(date string) {

	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Error("Failed to log to file, using default stderr")
	}

	log.SetLevel(logrus.InfoLevel)

	log.WithFields(logrus.Fields{
		"date": date,
	}).Trace("Init Log")
}
