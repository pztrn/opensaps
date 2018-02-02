// Flagger - arbitrary CLI flags parser.
//
// Copyright (c) 2017, Stanislav N. aka pztrn.
// All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package flagger

import (
    // stdlib
    "log"
    "os"
    "testing"
)

var (
    f *Flagger
)

func TestFlaggerInitialization(t *testing.T) {
    f = New(LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
    if f == nil {
        t.Fatal("Logger initialization failed!")
        t.FailNow()
    }
    f.Initialize()
}

func TestFlaggerAddBoolFlag(t *testing.T) {
    flag_testbool := Flag{
        Name: "testboolflag",
        Description: "Testing boolean flag",
        Type: "bool",
        DefaultValue: true,
    }
    err := f.AddFlag(&flag_testbool)
    if err != nil {
        t.Fatal("Failed to add boolean flag!")
        t.FailNow()
    }
}

func TestFlaggerAddIntFlag(t *testing.T) {
    flag_testint := Flag{
        Name: "testintflag",
        Description: "Testing integer flag",
        Type: "int",
        DefaultValue: 1,
    }
    err := f.AddFlag(&flag_testint)
    if err != nil {
        t.Fatal("Failed to add integer flag!")
        t.FailNow()
    }
}

func TestFlaggerAddStringFlag(t *testing.T) {
    flag_teststring := Flag{
        Name: "teststringflag",
        Description: "Testing string flag",
        Type: "string",
        DefaultValue: "superstring",
    }
    err := f.AddFlag(&flag_teststring)
    if err != nil {
        t.Fatal("Failed to add string flag!")
        t.FailNow()
    }
}

// This test doing nothing more but launching flags parsing.
func TestFlaggerParse(t *testing.T) {
    f.Parse()
}

func TestFlaggerGetBoolFlag(t *testing.T) {
    val, err := f.GetBoolValue("testboolflag")
    if err != nil {
        t.Fatal("Failed to get boolean flag: " + err.Error())
        t.FailNow()
    }

    if !val {
        t.Fatal("Failed to get boolean flag - should be true, but false received")
        t.FailNow()
    }
}

func TestFlaggerGetIntFlag(t *testing.T) {
    val, err := f.GetIntValue("testintflag")
    if err != nil {
        t.Fatal("Failed to get integer flag: " + err.Error())
        t.FailNow()
    }

    if val == 0 {
        t.Fatal("Failed to get integer flag - should be 1, but 0 received")
        t.FailNow()
    }
}

func TestFlaggerGetStringFlag(t *testing.T) {
    val, err := f.GetStringValue("teststringflag")
    if err != nil {
        t.Fatal("Failed to get string flag: " + err.Error())
        t.FailNow()
    }

    if val == "" {
        t.Fatal("Failed to get string flag - should be 'superstring', but nothing received")
        t.FailNow()
    }
}
