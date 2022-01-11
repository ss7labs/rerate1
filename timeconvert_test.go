package main
import (
 "testing"
)

func TestConvertTime (t *testing.T) {
 orgDateTime := convertTime("2020-06-30T14:35:03Z")
 expect := "2020-06-30"
 got := orgDateTime.Date
 if got != expect {
  t.Errorf("Expected %s, got %s", expect, got)
 }
 expect = "14:35"
 got = orgDateTime.Time
 if got != expect {
  t.Errorf("Expected %s, got %s", expect, got)
 }
}
