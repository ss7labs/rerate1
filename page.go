package main

import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func (org *Org) alignDstNumb(dst, numb string) string {
	l1 := len([]rune(dst))
	l2 := 24 - l1
	format := "%s%" + strconv.Itoa(l2) + "s"
	s := fmt.Sprintf(format, dst, numb)
	return s
}

func (org *Org) addLine(line string) {
	org.Page = append(org.Page, line)
}

func (org *Org) addHeader() {
	yy, _ := strconv.Atoi(org.Date[0:4])
	mm, _ := strconv.Atoi(org.Date[4:])

	t := time.Date(yy, time.Month(mm), 01, 00, 00, 00, 123456789, time.Now().Location())
	endOfMonth := now.With(t).EndOfMonth() // 2013-02-28 23:59:59.999999999 Thu
	dd := strconv.Itoa(endOfMonth.Day())

	orgid := strconv.Itoa(org.Id)
	orgacct := strconv.Itoa(org.Acct)

	org.Page = append(org.Page, "------------------------------------------------------------------------------")
	org.Page = append(org.Page, "  Лицевой счет :   "+orgacct+"             за пеpиод : 01."+org.Date[4:]+"."+org.Date[0:4]+"-"+dd+"."+org.Date[4:]+"."+org.Date[0:4])

	l := 77 - 2 - len([]rune(org.Name))
	format := "%2s%s%" + strconv.Itoa(l) + "s"
	s := fmt.Sprintf(format, "", org.Name, orgid)
	fmt.Println(s, l)
	org.Page = append(org.Page, s)
	//org.Page = append(org.Page,"  "+org.Name+"                                   "+orgid)

	org.Page = append(org.Page, "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -  - -- - - - ")
	org.Page = append(org.Page, "Телефон  Дата     Страна     Вызываемый       Время    К-во    Оплата   Оплата   ")
	org.Page = append(org.Page, "вызывающ.                     телефон        соедин.   мин.    (долл.)  (ман.)   ")
}

/*
func Output () {
 fmt.Println(len(org.Page))
 printorg.Page()
}
*/
func (org *Org) printPage() {
	for _, v := range org.Page {
		fmt.Println(v)
	}
}
func (org *Org) addFooter() bool {
	org.Page = append(org.Page, "    Всего по организации начислено :                    ")
	org.Page = append(org.Page, "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -  - -- - - - ")
	var s string
	format := "%4s%-54s%10.2f%10.2f" //    ТУРКМЕНИСТАН                                         "
	if org.IntCallsTotal.Man > 0 {
		s = fmt.Sprintf(format, "", "Междун.переговоры", org.IntCallsTotal.Usd, org.IntCallsTotal.Man)
		org.Page = append(org.Page, s)
	}
	if org.SngCallsTotal.Man > 0 {
		s = fmt.Sprintf(format, "", "СНГ", org.SngCallsTotal.Usd, org.SngCallsTotal.Man)
		org.Page = append(org.Page, s)
	}
	if org.TkmCallsTotal.Man > 0 {
		s = fmt.Sprintf(format, "", "ТУРКМЕНИСТАН", org.TkmCallsTotal.Usd, org.TkmCallsTotal.Man)
		org.Page = append(org.Page, s)
	}
	//org.Page = append(org.Page,"    ТУРКМЕНИСТАН                                              114.20     75.18")
	org.Page = append(org.Page, "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -  - -- - - - ")

	totalCallsUsd := org.SngCallsTotal.Usd + org.IntCallsTotal.Usd + org.TkmCallsTotal.Usd
	totalCallsMan := +org.SngCallsTotal.Man + org.IntCallsTotal.Man + org.TkmCallsTotal.Man
	if totalCallsMan > 0 {
		s = fmt.Sprintf(format, "", "Всего за меж/гор. меж/нар. переговоры в кpедит :", totalCallsUsd, totalCallsMan)
		org.Page = append(org.Page, s)
	}
	//org.Page = append(org.Page,"    Всего за меж/гор. меж/нар. переговоры в кpедит :          114.20     75.18")

	org.tenPercentUsd = 0.0
	org.tenPercentMan = 0.0

	if org.Firma != 3 {
		org.tenPercentUsd = totalCallsUsd * 0.1
		org.tenPercentMan = totalCallsMan * 0.1
	}

	if org.tenPercentMan > 0 {
		s = fmt.Sprintf(format, "", "10%  за меж/гор. меж/нар. переговоры в кpедит :", org.tenPercentUsd, org.tenPercentMan)
		org.Page = append(org.Page, s)
	}
	//org.Page = append(org.Page,"    10%  за меж/гор. меж/нар. переговоры в кpедит :            11.42      7.52")

	if org.LocalCallsTotal.Man > 0 {
		s = fmt.Sprintf(format, "", "Cверхлимитные переговоры по городу :", org.LocalCallsTotal.Usd, org.LocalCallsTotal.Man)
		org.Page = append(org.Page, s)
	}
	//org.Page = append(org.Page,"    Cверхлимитные переговоры по городу :                        0.36      1.05")

	if org.PstnServiceTotal.Man > 0 {
		s = fmt.Sprintf(format, "", "Услуги цифровых АТС", org.PstnServiceTotal.Usd, org.PstnServiceTotal.Man)
		org.Page = append(org.Page, s)
	}
	//org.Page = append(org.Page,"    Услуги цифровых АТС                                         0.00     32.00")

	totalUsd := totalCallsUsd + org.tenPercentUsd + org.LocalCallsTotal.Usd + org.PstnServiceTotal.Usd
	totalMan := totalCallsMan + org.tenPercentMan + org.LocalCallsTotal.Man + org.PstnServiceTotal.Man

	if totalMan <= 0 {
		return false
	}

	length := len(org.Page) + 3
	listCount := length / 63
	if length%63 > 0 {
		listCount++
	}
	if listCount == 0 {
		listCount++
	}

	listPriceMan := float64(listCount * 2)
	listPriceUsd := listPriceMan / 3.5

	log.WithFields(logrus.Fields{
		"length":    length,
		"listCount": listCount,
	}).Debug("addFooter")

	formatCount := "%s%6d.    Плата за pасшифpовку :"
	countStr := fmt.Sprintf(formatCount, "Печатных листов :", listCount)
	s = fmt.Sprintf(format, "", countStr, listPriceUsd, listPriceMan)
	org.Page = append(org.Page, s)
	//org.Page = append(org.Page,"    Печатных листов :     5.    Плата за pасшифpовку :          2.85     10.00")

	org.orgTotalUsd = totalUsd
	org.orgTotalMan = totalMan
	totalUsd += listPriceUsd
	totalMan += listPriceMan
	s = fmt.Sprintf(format, "", "Сумма к оплате  :", totalUsd, totalMan)
	org.Page = append(org.Page, s)
	//org.Page = append(org.Page,"    Сумма к оплате  :                                         128.83    125.75")

	org.Page = append(org.Page, "------------------------------------------------------------------------------")
	return true
}
