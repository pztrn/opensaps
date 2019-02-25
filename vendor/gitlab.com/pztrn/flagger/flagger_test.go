// Flagger - arbitrary CLI flags parser.
//
// Copyright (c) 2017-2018, Stanislav N. aka pztrn.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject
// to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package flagger

import (
	// stdlib
	"log"
	"os"
	"testing"

	// other
	"github.com/stretchr/testify/require"
)

func TestFlaggerInitialization(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()
}

func TestFlaggerInitializationWithNilLogger(t *testing.T) {
	f := New("tests", nil)
	require.NotNil(t, f)
	f.Initialize()
}

func TestFlaggerAddBoolFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestBool := Flag{
		Name:         "testboolflag",
		Description:  "Testing boolean flag",
		Type:         "bool",
		DefaultValue: true,
	}
	err := f.AddFlag(&flagTestBool)
	require.Nil(t, err)
}

func TestFlaggerAddSameBoolVar(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestBool := Flag{
		Name:         "testboolflag",
		Description:  "Testing boolean flag",
		Type:         "bool",
		DefaultValue: true,
	}
	err := f.AddFlag(&flagTestBool)
	require.Nil(t, err)

	flagTestBool1 := Flag{
		Name:         "testboolflag",
		Description:  "Testing boolean flag",
		Type:         "bool",
		DefaultValue: true,
	}
	err1 := f.AddFlag(&flagTestBool1)
	require.NotNil(t, err1)
}

func TestFlaggerAddIntFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestInt := Flag{
		Name:         "testintflag",
		Description:  "Testing integer flag",
		Type:         "int",
		DefaultValue: 1,
	}
	err := f.AddFlag(&flagTestInt)
	require.Nil(t, err)
}

func TestFlaggerAddStringFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestString := Flag{
		Name:         "teststringflag",
		Description:  "Testing string flag",
		Type:         "string",
		DefaultValue: "superstring",
	}
	err := f.AddFlag(&flagTestString)
	require.Nil(t, err)
}
func TestFlaggerParse(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestString := Flag{
		Name:         "teststringflag",
		Description:  "Testing string flag",
		Type:         "string",
		DefaultValue: "superstring",
	}
	err := f.AddFlag(&flagTestString)
	require.Nil(t, err)

	f.Parse()
}

func TestFlaggerParseAndReparse(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestString := Flag{
		Name:         "teststringflag",
		Description:  "Testing string flag",
		Type:         "string",
		DefaultValue: "superstring",
	}
	err := f.AddFlag(&flagTestString)
	require.Nil(t, err)

	f.Parse()
	f.Parse()
}

func TestFlaggerGetBoolFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestBool := Flag{
		Name:         "testboolflag",
		Description:  "Testing boolean flag",
		Type:         "bool",
		DefaultValue: true,
	}
	err := f.AddFlag(&flagTestBool)
	require.Nil(t, err)

	f.Parse()

	val, err := f.GetBoolValue("testboolflag")
	require.Nil(t, err)
	require.True(t, val)
}

func TestFlaggerGetUnknownBoolFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestBool := Flag{
		Name:         "testboolflag",
		Description:  "Testing boolean flag",
		Type:         "bool",
		DefaultValue: true,
	}
	err := f.AddFlag(&flagTestBool)
	require.Nil(t, err)

	f.Parse()

	val, err := f.GetBoolValue("unknownboolflag")
	require.NotNil(t, err)
	require.False(t, val)
}

func TestFlaggerGetIntFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestInt := Flag{
		Name:         "testintflag",
		Description:  "Testing integer flag",
		Type:         "int",
		DefaultValue: 1,
	}
	err := f.AddFlag(&flagTestInt)
	require.Nil(t, err)

	f.Parse()

	val, err := f.GetIntValue("testintflag")
	require.Nil(t, err)
	require.NotEqual(t, 0, val)
}

func TestFlaggerGetUnknownIntFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestInt := Flag{
		Name:         "testintflag",
		Description:  "Testing integer flag",
		Type:         "int",
		DefaultValue: 1,
	}
	err := f.AddFlag(&flagTestInt)
	require.Nil(t, err)

	f.Parse()

	val, err := f.GetIntValue("unknownintflag")
	require.NotNil(t, err)
	require.Equal(t, 0, val)
}

func TestFlaggerGetStringFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestString := Flag{
		Name:         "teststringflag",
		Description:  "Testing string flag",
		Type:         "string",
		DefaultValue: "superstring",
	}
	err := f.AddFlag(&flagTestString)
	require.Nil(t, err)

	f.Parse()

	val, err := f.GetStringValue("teststringflag")
	require.Nil(t, err)
	require.NotEqual(t, "", val)
	require.Equal(t, "superstring", val)
}

func TestFlaggerGetUnknownStringFlag(t *testing.T) {
	f := New("tests", LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
	require.NotNil(t, f)
	f.Initialize()

	flagTestString := Flag{
		Name:         "teststringflag",
		Description:  "Testing string flag",
		Type:         "string",
		DefaultValue: "superstring",
	}
	err := f.AddFlag(&flagTestString)
	require.Nil(t, err)

	f.Parse()

	val, err := f.GetStringValue("unknownstringflag")
	require.NotNil(t, err)
	require.Equal(t, "", val)
}
