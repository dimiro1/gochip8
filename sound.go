// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

// The sound engine interface
// The concrete implementation must implement this interface
type sound interface {
	Beep()
}
