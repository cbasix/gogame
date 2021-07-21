// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type Coordinator interface {
	register() chan<- PlayerViewListener
	unregister() chan<- PlayerViewListener
}
