package main

import (
	"regexp"
	//"fmt"
)

type OrgDetailTime struct {
	Date string
	Time string
}

func convertTime(datetime string) OrgDetailTime {
	odt := OrgDetailTime{}
	r := regexp.MustCompile(`(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})T(?P<HH>\d{2}):(?P<MM>\d{2}):(?P<SS>\d{2})Z`)
	s := r.FindStringSubmatch(datetime)
	odt.Date = s[1] + "-" + s[2] + "-" + s[3]
	odt.Time = s[4] + ":" + s[5]
	//fmt.Printf("%#v\n", r.FindStringSubmatch(datetime))
	//fmt.Printf("%#v\n", r.SubexpNames())
	return odt
}
