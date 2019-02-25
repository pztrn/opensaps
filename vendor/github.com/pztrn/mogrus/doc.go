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

// Mogrus is a logger package built on top of Logrus with extended
// (and kinda simplified) approach for multilogging.
//
// Example usage:
//
//     TBF
//
// Output functions are splitted in three things:
//   * ``Debug`` - simple debug printing
//   * ``Debugf`` - debug printing with line formatting
//   * ``Debugln`` - debug printing, same as Debug().
//
// It will try to resemble Logrus as much as possible.
//
// About functions.
//
// Info(f,ln) and Print(f,ln) doing same thing - logging to
// INFO log level.
//
// Be careful while using Fatal(f,ln) functions, because
// it will log only to first logger and then will os.Exit(1)!
// There is no guarantee what logger will be first ATM.
