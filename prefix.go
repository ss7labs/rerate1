package main

import (
	"database/sql"
	"regexp"
	"strings"
	//"fmt"
	"github.com/sirupsen/logrus"
)

var orderedPrefix []string

type SplittedNumber struct {
	Direction string
	Remained  string
}
type SearchDirection interface {
	DirectionByPrefix(prefix string) string
}

type NumberDirection string

func (nd NumberDirection) DirectionByPrefix(prefix string) string {
	return getNameOfPrefix(prefix)
}

func SplitNumber(number string, sd SearchDirection) SplittedNumber {
	sn := SplittedNumber{}
	var prefix string
	m, _ := regexp.MatchString(`^[1-9]\d{5}$`, number)
	if m {
		prefix = "80012"
		sn.Remained = number
	} else {
		prefix = matchedPrefix(number)
		s := strings.Split(number, prefix)
		sn.Remained = s[1]
	}
	sn.Direction = sd.DirectionByPrefix(prefix)
	return sn
}

func matchedPrefix(number string) string {
	prefix := ""
	for _, value := range orderedPrefix {
		regex := "^" + value + "[0-9]+"
		matched, _ := regexp.MatchString(regex, number)
		if matched {
			prefix = value
			break
		}
	}
	return prefix
}
func getDirection(number string) string {
	direction := ""

	for _, value := range orderedPrefix {
		regex := "^" + value + "[0-9]+"
		matched, _ := regexp.MatchString(regex, number)
		if matched {
			direction = getNameOfPrefix(value)
			break
		}
	}

	return direction
}

func getNameOfPrefix(prefix string) string {
	output := ""
	shortname := ""
	query := "SELECT name,shortname FROM route_prices WHERE prefix='" + prefix + "'"
	err := asudb.QueryRow(query).Scan(&output, &shortname)
	if sqlErrCheck(err) {
		output = ""
	}
	length := len([]rune(output))

	log.WithFields(logrus.Fields{
		"output": output,
		"length": length,
	}).Debug("getNameOfPrefix")

	if length > 12 {
		output = shortname
	}
	return output
}

func sqlErrCheck(err error) bool {
	norows := false
	if err != nil {
		if err == sql.ErrNoRows {
			norows = true
		} else {
			panic(err.Error())
		}
	}
	return norows
}

func loadPrefixes() {
	query := "SELECT prefix FROM route_prices ORDER BY CHAR_LENGTH(prefix) DESC"

	rows, err := asudb.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var prefix string
		if err := rows.Scan(&prefix); err != nil {
			panic(err.Error())
		}
		orderedPrefix = append(orderedPrefix, prefix)
	}
}
