package main

import (
	"fmt"
	"strconv"
)

type Dual struct {
	totalUsd float64
	totalMan float64
	usdPage  []string
	manPage  []string
	usdFn    string
	manFn    string
}

func (org *Org) printDualOrg() {
	fmt.Println("printDualOrg", org.Id, org.Name, org.Acct)
	org.dualOrg = Dual{}
	org.dualIntCalls()
	org.dualLocalCalls()
}

func (org *Org) dualIntCalls() {
	inBlock := org.dualPhoneListINSQL()
	fmt.Println(org.Id, inBlock)
	query := "SELECT numa,sum(rate_usd) FROM rated WHERE toYYYYMM(event_date)='" + org.Date + "' AND numb LIKE '810%' AND numa IN " + inBlock + " GROUP BY numa ORDER BY numa ASC"

	rows, err := chdb.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var numa string
		var usd float64
		if err := rows.Scan(&numa, &usd); err != nil {
			panic(err.Error())
		}
		if usd == 0 {
			continue
		}
		org.dualIntPrintDetail(numa, usd)
	}
	if org.dualOrg.totalUsd <= 0 {
		return
	}
	//org.dualSaveFile()
}

func (org *Org) dualLocalCalls() {
}

func (org *Org) dualIntPrintDetail(numa string, usd float64) {
	/*
	   var str string
	   if org.isFirstLine() {
	    str = numa
	   }
	*/
	fmt.Println("dualIntPrint", numa, usd)
	org.dualOrg.totalUsd += usd
	/* query := "SELECT toDateTime(event_time,'UTC'),numb,duration,rate_usd,rate_org FROM rated WHERE toYYYYMM(event_date)='"+org.Date+"' AND numb LIKE '810%' AND numb NOT LIKE '8107%' AND numa='"+phone+"' ORDER BY event_time"
	   org.callsDetails(query,phone)
	   s := fmt.Sprintf("%7s%-51s%10.2f%10.2f", "","Всего по телефону меж.переговоров",totalUsd,totalOrg)
	   org.addLine(s)
	*/
}

func (org *Org) dualPhoneListINSQL() string {
	query := "SELECT phone FROM phones WHERE org_id=" + strconv.Itoa(org.Id)
	rows, err := asudb.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	str := "('"
	for rows.Next() {
		var numa string
		if err := rows.Scan(&numa); err != nil {
			panic(err.Error())
		}
		str = str + numa + "','"
	}
	runes := []rune(str)
	length := len(runes)
	str = string(runes[:length-3])
	str += "')"
	return str
}
