package main
import (
 "fmt"
 "strconv"
 "math"
 "database/sql"
 _ "github.com/go-sql-driver/mysql"
 _ "github.com/mailru/go-clickhouse"
)
type CallDetail struct {
 DTime string
 Numa string
 Numb string
 DurSec int
 Usd float64
 Man float64
 PrintNuma bool
}
type TotalDual struct {
 Usd float64
 Man float64
}
type PhonePrintControl struct {
 DetailsBlockExists bool
 LocalsBlockExists bool
 PstnExists bool
}
type Org struct {
 lines []string
 dualOrg Dual
 Id int
 Name string
 Acct int
 Date string
 Firma int
 IntCallsTotal TotalDual
 SngCallsTotal TotalDual
 TkmCallsTotal TotalDual
 LocalCallsTotal TotalDual
 PstnServiceTotal TotalDual
 Page []string
 phoneTotalUsd float64
 phoneTotalOrg float64
 tenPercentUsd float64
 tenPercentMan float64
 orgTotalUsd float64
 orgTotalMan float64
 phonePrintControl *PhonePrintControl
 dual bool
 reestr *Reestr
}

func (org *Org) printOneOrg() {
 fmt.Println(org.Id, org.Name, org.Acct)
 query := "SELECT phone FROM phones WHERE org_id="+strconv.Itoa(org.Id)+" ORDER BY phone ASC"

 rows, err := asudb.Query(query)                                                                      
 if err != nil {
  panic(err.Error())
 }                    
 defer rows.Close()

 org.addHeader()                                                  
 for rows.Next() {
  var phone string
  if err := rows.Scan(&phone); err != nil {
    panic(err.Error())
  }
  org.printOnePhone(phone)
 } 
 if org.addFooter()  {                                                 
  org.printPageToFile()
  org.addToReestr()
 }
 //org.printPage()
}

func (org *Org) printOnePhone(phone string) {
 org.phoneTotalUsd = 0
 org.phoneTotalOrg = 0
 org.phonePrintControl = &PhonePrintControl{DetailsBlockExists:false,LocalsBlockExists:false,PstnExists:false}
/*
** Get Total of 810 calls and NOT 8107
*/
 query := "SELECT sum(rate_usd),sum(rate_org) FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numb LIKE '810%' AND numb NOT LIKE '8107%' AND numa='"+phone+"' GROUP BY numa"
 totalUsd,totalOrg := getSum(query)

 if totalOrg > 0 {
  //fmt.Println("International total",phone,totalUsd,totalOrg)
  org.IntCallsTotal.Usd += totalUsd
  org.IntCallsTotal.Man += totalOrg
  org.phoneTotalUsd += totalUsd
  org.phoneTotalOrg += totalOrg
  org.phonePrintControl.DetailsBlockExists = true
  query := "SELECT event_time,numb,duration,rate_usd,rate_org FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numb LIKE '810%' AND numb NOT LIKE '8107%' AND numa='"+phone+"' ORDER BY event_time"
  org.callsDetails(query,phone)
  s := fmt.Sprintf("%7s%-51s%10.2f%10.2f", "","Всего по телефону меж.переговоров",totalUsd,totalOrg) 
  org.addLine(s)
 }

/*
** Get Total of 8107 calls
*/
 query = "SELECT sum(rate_usd),sum(rate_org) FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numb LIKE '8107%' AND numa='"+phone+"' GROUP BY numa"
 totalUsd,totalOrg = getSum(query)

 if totalOrg > 0 {
  fmt.Println("SNG total",phone,totalUsd,totalOrg)
  org.SngCallsTotal.Usd += totalUsd
  org.SngCallsTotal.Man += totalOrg
  org.phoneTotalUsd += totalUsd
  org.phoneTotalOrg += totalOrg
  org.phonePrintControl.DetailsBlockExists = true
  query := "SELECT event_time,numb,duration,rate_usd,rate_org FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numb LIKE '8107%' AND numa='"+phone+"' ORDER BY event_time"
  org.callsDetails(query,phone)
  s := fmt.Sprintf("%7s%-51s%10.2f%10.2f", "","Всего по телефону разговоров по СНГ",totalUsd,totalOrg) 
  org.addLine(s)
 }

/*
** Get Total of non 810 and 8 calls
*/
 query = "SELECT sum(rate_usd),sum(rate_org) FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numb LIKE '8%' AND numb NOT LIKE '810%' AND numa='"+phone+"' GROUP BY numa"
 totalUsd,totalOrg = getSum(query)

 if totalOrg > 0 {
  //fmt.Println("Turkmenistan total",phone,totalUsd,totalOrg)
  org.TkmCallsTotal.Usd += totalUsd
  org.TkmCallsTotal.Man += totalOrg
  org.phoneTotalUsd += totalUsd
  org.phoneTotalOrg += totalOrg
  org.phonePrintControl.DetailsBlockExists = true
  query := "SELECT event_time,numb,duration,rate_usd,rate_org FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numb LIKE '8%' AND numb NOT LIKE '810%' AND numa='"+phone+"' ORDER BY event_time"
  org.callsDetails(query,phone)
  s := fmt.Sprintf("%7s%-51s%10.2f%10.2f", "","Всего по телефону разговоров по Туркменистану",totalUsd,totalOrg) 
  org.addLine(s)
 }
 //fmt.Println("Turkmenistan total",phone,org.TkmCallsTotal.Usd,org.TkmCallsTotal.Man)

/*
** Get Total of Local non 810 and non 8 calls
*/
 query = "SELECT sum(rate_usd),sum(rate_org) FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numb NOT LIKE '8%' AND numa='"+phone+"' GROUP BY numa"
 totalUsd,totalOrg = 0.0,0.0

 if org.Firma != 3 {
  totalUsd,totalOrg = getSum(query)
 }

 if totalOrg > 0 {
  intP, frac := math.Modf(totalOrg*100)

  //fmt.Println(intP,frac)

  if frac > 0 { intP++ }
  totalOrg = float64(intP/100)

  if totalUsd == 0 {
   totalUsd = convertToUsd(totalOrg)
  }

//Total for whole Org
  org.LocalCallsTotal.Usd += totalUsd
  org.LocalCallsTotal.Man += totalOrg
//Total for one phone
  org.phoneTotalUsd += totalUsd
  org.phoneTotalOrg += totalOrg
  org.phonePrintControl.LocalsBlockExists = true

  query = "SELECT sum(exc_dur) FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numa='"+phone+"' GROUP BY numa"
  localSec := getLocalTimeSum(query)
  var localMin int
  if (localSec % 60) > 0 {
   localMin = localSec/60 + 1
  } else { localMin = localSec/60 }
  s := fmt.Sprintf("%s Свеpхлимитные переговоры по гоpоду%17d%10.2f%10.2f", phone,localMin,totalUsd,totalOrg) 

  if org.phonePrintControl.DetailsBlockExists {
   s = fmt.Sprintf("%7sСвеpхлимитные переговоры по гоpоду%17d%10.2f%10.2f", "",localMin,totalUsd,totalOrg) 
  }

  //fmt.Println("Свеpхлимитные переговоры по гоpоду",phone,localMin,totalUsd,totalOrg)
  org.addLine(s)
 }
/*
** Get PSTN Services usage
** pstnservices.go
*/
 org.addPstnServices(phone)

/*
** add Total lines of One phone
*/
 if org.phonePrintControl.DetailsBlockExists || org.phonePrintControl.LocalsBlockExists || org.phonePrintControl.PstnExists {
  s := fmt.Sprintf("%7s%-51s%10.2f%10.2f\n", "","Всего по телефону :",org.phoneTotalUsd,org.phoneTotalOrg) 
  org.addLine(s)
 }
}

func (org *Org) callsDetails(query, phone string) {
 rows, err := chdb.Query(query)                                                                      
 if err != nil {
  panic(err.Error())
 }                    
 defer rows.Close()
 count := 0                                                  
 for rows.Next() {
  var cd CallDetail
  if err := rows.Scan(&cd.DTime,&cd.Numb,&cd.DurSec,&cd.Usd,&cd.Man); err != nil {
    panic(err.Error())
  }
  if count == 0 {
   cd.PrintNuma = true
   cd.Numa = phone
   count++
  }
  org.prepareDetailRow(cd)
 } 
}
func (org *Org) prepareDetailRow(cd CallDetail) {
 sd  := NumberDirection(cd.Numb)
 spn := SplitNumber(cd.Numb,sd)
 odt := convertTime(cd.DTime)

 //fmt.Println(odt.Date,spn.Direction,spn.Remained,odt.Time,secToMin(cd.DurSec),cd.Usd,cd.Man)
 dateRunes := []rune(odt.Date) //2020-06-25
 yyyy, _ := strconv.Atoi(string(odt.Date[0:4]))
 yy := yyyy-2000

 //Convert to 25.06.20
 date := string(dateRunes[8:])+"."+string(dateRunes[5:7])+"."+strconv.Itoa(yy)
 alignedDstNumb := org.alignDstNumb(spn.Direction,spn.Remained)
 format := "%-7s%s %s%12s%6d%10.2f%10.2f"
 s := fmt.Sprintf(format, "",date,alignedDstNumb,odt.Time,secToMin(cd.DurSec),cd.Usd,cd.Man) 
 if cd.PrintNuma {
//950589  2.06.20 Ашх.(сот.)       1746266       12: 3     2      0.40      0.28
  s = fmt.Sprintf(format,cd.Numa,date,alignedDstNumb,odt.Time,secToMin(cd.DurSec),cd.Usd,cd.Man) 
  cd.PrintNuma = false
 }
 org.addLine(s)
}

func getLocalTimeSum (query string) (totalMin int) {

 err := chdb.QueryRow(query).Scan(&totalMin)

 if err != nil {
  if err == sql.ErrNoRows {
   totalMin=0
  } else {
   panic(err.Error())
  }
 }
 return totalMin
}

func getSum (query string) (totalUsd,totalOrg float64) {

 err := chdb.QueryRow(query).Scan(&totalUsd,&totalOrg)

 if err != nil {
  if err == sql.ErrNoRows {
   totalUsd=0
   totalOrg=0
  } else {
   panic(err.Error())
  }
 }
 return totalUsd,totalOrg
}

func convertToUsd (man float64) (usd float64) {
 if man == 0 { 
  usd = 0 
  return
 }
 intP,_ := math.Modf((man/2.85)*100)
 usd = intP/100
 return
}
