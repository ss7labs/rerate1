package main

import (
 "testing"
 _"fmt"
)

type mockND struct {}

func (mnd mockND) DirectionByPrefix(prefix string) string {
 return "Ашх.(сот.)"
}

func TestSplitNumber(t *testing.T) {
 var sn SplittedNumber
 orderedPrefix = append(orderedPrefix,"8006")

 mock := &mockND{}
 sn = SplitNumber("80062606182",mock)

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
}