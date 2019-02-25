package main

import (
	"gitlab.com/pztrn/flagger"
)

var f *flagger.Flagger

func main() {
	f = flagger.New("testprogram", nil)
	f.Initialize()
	f.AddFlag(&flagger.Flag{
		Name:         "testflag",
		Description:  "Just a test flag",
		Type:         "bool",
		DefaultValue: false,
	})
	f.Parse()
}
