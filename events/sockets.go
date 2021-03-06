// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package events

import "github.com/gorilla/websocket"

// WSManager wraps a websocket connection and should be managed inside of a
// single go routine. Websockets are not safe for Concurrent reads and writes
// and so the WSManager is used to guarantee that only one go routine is
// reading or writing at any given time.
type WSManager struct {
	In     chan []byte
	Out    chan []byte
	Socket *websocket.Conn
	InBuf  [1024]byte
	OutBuf [1024]byte
}
