package main

import (   
 "math"    
)

func round(in float64) float64 {
 out := 0.0
 fract0 := math.Trunc(in)
 fract1 := int(math.Trunc((in - math.Trunc(in))*1000))
 rem1 := fract1 % 100
 rem2 := rem1 % 10
 rem := fract1-rem1+rem1-rem2
 if rem2 > 0 {
  rem += 10
 }
 out = fract0+float64(rem)/1000
 return out
}

func secToMin(sec int) int {
 min := 0
 if sec == 0 {
  return min
 }

 if sec <= 60 {
  min = 1
  return min
 }

 if sec > 60 {
  remainder := sec%60
  min = (sec-remainder)/60
  if remainder > 0 {
   min += 1
  }
 }
 return min
}
