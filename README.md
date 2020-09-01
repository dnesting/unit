# unit

[![GoDoc](https://godoc.org/github.com/dnesting/unit?status.svg)](https://godoc.org/github.com/dnesting/unit)
[![Build Status](https://travis-ci.org/dnesting/unit.svg?branch=master)](https://travis-ci.org/dnesting/unit)
[![codecov](https://coveralls.io/repos/github/dnesting/unit/badge.svg?branch=master)](https://coveralls.io/github/dnesting/unit?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dnesting/unit)](https://goreportcard.com/report/github.com/dnesting/unit)

This is a work in progress.

The goal is to enable light-weight operations on values with units attached to them, including implicit unit conversion.

For instance, the following should be possible:

```go
mps := si.Meter.Div(si.Second)
speed := mps(5)  // 5 m/s
dist := speed.Mul(si.FromDuration(5*time.Second))
fmt.Println(dist)  // "25 m"

mph := us.Mile.Div(us.Hour)
speed = speed.Convert(mph) // now "11.185 mi/hr"
fmt.Println(dist.Div(speed))  // "5 s"
```

Stretch goal is compatibility with GNU Units `definitions.units`, something like:

```go
defs := unitdef.Definitions()  // read and parse definitions.units
mps = defs.Must("m/s")
dist := mps(5).Mul(si.Second(5))
fmt.Println(dist)  // "25 m"
```
