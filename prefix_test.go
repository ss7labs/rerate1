package main

import (
	_ "fmt"
	"testing"
)

type mockND struct{}

func (mnd mockND) DirectionByPrefix(prefix string) string {
	codeName := "Ашх.(сот.)"
	if prefix == "80012" {
		codeName = "Ашхабад"
	}
	return codeName
}

func TestSplitNumber(t *testing.T) {
	var sn SplittedNumber
	orderedPrefix = append(orderedPrefix, "8006")
	orderedPrefix = append(orderedPrefix, "80012")

	mock := &mockND{}
	sn = SplitNumber("80062606182", mock)

	expect := "2606182"
	got := sn.Remained
	if got != expect {
		t.Errorf("Remained expected %s, got %s", expect, got)
	}

	expect = "Ашх.(сот.)"
	got = sn.Direction
	if got != expect {
		t.Errorf("Direction expected %s, got %s", expect, got)
	}

	sn = SplitNumber("414747", mock)

	expect = "414747"
	got = sn.Remained
	if got != expect {
		t.Errorf("Remained expected %s, got %s", expect, got)
	}

	expect = "Ашхабад"
	got = sn.Direction
	if got != expect {
		t.Errorf("Direction expected %s, got %s", expect, got)
	}
}
