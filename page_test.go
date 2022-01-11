package main

import (
	"fmt"
	"testing"
)

func createOrg() *Org {
	org := &Org{}
	org.Acct = 1
	org.Id = 1288
	org.Name = "Китайская Национальная Нефтегазовая Корпорация."
	org.Date = "202006"
	return org
}
func TestPageHeader(t *testing.T) {
	org := createOrg()
	org.addHeader()
	dst := "Ашх.(сот.)"
	fmt.Println(dst, len(dst))
	alignedDstNumb := org.alignDstNumb(dst, "1746266")
	s := fmt.Sprintf("%-7s%s %s%12s%6d%10.2f%10.2f", "950589", "02.06.20", alignedDstNumb, "17:03", 2, 0.40, 0.28)
	org.addLine(s)
	s = fmt.Sprintf("%-7s%s %s%14s%12s%6d%10.2f%10.2f", "", "02.06.20", "Ашх.(сот.)", "1746266", "17:03", 2, 0.40, 0.28)
	org.addLine(s)

	s = fmt.Sprintf("%s Свеpхлимитные переговоры по гоpоду%17d%10.2f%10.2f", "233127", 8, 0.00, 0.01)
	org.addLine(s)
	s = fmt.Sprintf("%s Свеpхлимитные переговоры по гоpоду%17d%10.2f%10.2f", "950582", 159, 0.03, 0.10)
	org.addLine(s)

	org.printPage()
}

func TestPageFooter(t *testing.T) {
	org := createOrg()
	org.SngCallsTotal.Usd = 0
	org.SngCallsTotal.Man = 0
	org.IntCallsTotal.Usd = 0
	org.IntCallsTotal.Man = 0
	org.TkmCallsTotal.Usd = 114.20
	org.TkmCallsTotal.Man = 75.18
	org.LocalCallsTotal.Usd = 0.36
	org.LocalCallsTotal.Man = 1.05

	org.PstnServiceTotal.Usd = 0
	org.PstnServiceTotal.Man = 32
	org.addFooter()
	org.printPage()
}

func TestAlignDstNumb(t *testing.T) {
	org := &Org{}
	expect := 24
	got := len([]rune(org.alignDstNumb("Ашх.(сот.)", "1746266")))
	if got != expect {
		t.Errorf("Aligned string length  expected %d, got %d", expect, got)
	}

	got = len([]rune(org.alignDstNumb("Ашх.(сот.)", "21746266")))
	if got != expect {
		t.Errorf("Aligned string length  expected %d, got %d", expect, got)
	}

}

func TestFakeHeader(t *testing.T) {
}
