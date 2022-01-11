package main

import (
 "fmt"
 "github.com/gammazero/workerpool"
  "database/sql"
 _ "github.com/go-sql-driver/mysql"
 _ "github.com/mailru/go-clickhouse"
 "github.com/sirupsen/logrus"
 "os"
)

var asudb *sql.DB
var chdb *sql.DB
var log = logrus.New()

func main() {

 if len(os.Args) < 2 {
  fmt.Printf("Usage: %s <YYYYMM> <dual>(optional)\n", os.Args[0])
  return
 } 

 date := os.Args[1]

 var dual string
 if len(os.Args) == 3 {
  dual = os.Args[2]
 }

 initLog(date)
 initDB()
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
   runWork(r.orgId,r.date,r.dual,reestr)
  })
 }

 wp.StopWait()
 reestr.closeFiles()
}

func initDB() {
 var err error
//Init clickhouse
 chdb, err = sql.Open("clickhouse", "http://default:te410pte412p@10.19.64.124:8123/asubill?debug=true")
 if err != nil {
  panic(err.Error())
 }
//Init ASUDB
 asudb, err = sql.Open("mysql", "asu:qwerty@tcp(localhost:3308)/agts_asu")
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
