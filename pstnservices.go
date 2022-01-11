package main

import (
	"database/sql"
	"fmt"
)

func (org *Org) addPstnServices(phone string) {
	query := "SELECT IFNULL(sum(spr.tarif2),0),IFNULL(sum(spr.tarif_usd),0) FROM `dvo_srv` srv INNER JOIN dvo_spr spr ON spr.dvo_id=srv.dvo_id WHERE phone='" + phone + "'"
	pstnAmount, pstnAmountUsd := getPstnSrvPrice(query)
	if pstnAmount <= 0 {
		return
	}

	s := fmt.Sprintf("%s Услуги ATC%51.2f%10.2f", phone, pstnAmountUsd, pstnAmount)

	if org.phonePrintControl.DetailsBlockExists || org.phonePrintControl.LocalsBlockExists {
		s = fmt.Sprintf("%7sУслуги ATC%51.2f%10.2f", "", pstnAmountUsd, pstnAmount)
	}

	org.phoneTotalUsd += pstnAmountUsd
	org.PstnServiceTotal.Usd += pstnAmountUsd

	org.phoneTotalOrg += pstnAmount
	org.PstnServiceTotal.Man += pstnAmount
	org.phonePrintControl.PstnExists = true
	org.addLine(s)
}

func getPstnSrvPrice(query string) (amount, amountUsd float64) {
	err := asudb.QueryRow(query).Scan(&amount, &amountUsd)

	if err != nil {
		if err == sql.ErrNoRows {
			amount = 0
		} else {
			panic(err.Error())
		}
	}
	return
}
