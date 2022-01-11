package main
import (
 "os"
 "strconv"
 "github.com/sirupsen/logrus"
)
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func (org *Org) printPageToFile () {
 runes := []rune(org.Date)
 month := string(runes[4:])
 orgacct := strconv.Itoa(org.Acct)[0:6]
 ext:= "R"
 if org.Firma == 3 {ext="B"}
 if org.Firma == 5 {ext="D"}
 fn := "L_"+orgacct+"."+ext+month

 log.WithFields(logrus.Fields{
      "fn": fn,
      "Acct": org.Acct,
      "orgacct":   orgacct,
 }).Debug("printPageToFile")

 f, err := os.Create(fn)
 check(err)
 defer f.Close()

 for _, v := range org.Page {
  line := v+"\n"
  f.WriteString(line)
 }
}
