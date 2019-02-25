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
package mogrus

type MogrusLogger struct {
	// Initialized loggers.
	// Key is a name of logger.
	loggers map[string]*LoggerHandler
}

// CreateLogger creates new logger handler, adds it to list of known
// loggers and return it to caller.
// Note that logger handler will be "just initialized", to actually
// use it you should add output with LoggerHandler.CreateOutput().
func (ml *MogrusLogger) CreateLogger(name string) *LoggerHandler {
	lh := &LoggerHandler{}
	lh.Initialize()
	ml.loggers[name] = lh

	return lh
}

// Initialize initializes Mogrus instance.
func (ml *MogrusLogger) Initialize() {
	ml.loggers = make(map[string]*LoggerHandler)
}
