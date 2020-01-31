// Copyright (C) 2019 LEAP
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

package sip2

import (
	"github.com/reiver/go-telnet"
)

// The terminator can be configured differently for different SIP endpoints.
// This gets set in sip2.auth according to an environment variable

var telnetTerminator string

func telnetRead(conn *telnet.Conn) (out string) {
	var buffer [1]byte
	recvData := buffer[:]
	var n int
	var err error

	for {
		n, err = conn.Read(recvData)
		if n <= 0 || err != nil {
			break
		} else {
			out += string(recvData)
		}
		if len(out) > 1 && out[len(out)-len(telnetTerminator):] == telnetTerminator {
			break
		}
	}
	return out
}

func telnetSend(conn *telnet.Conn, command string) {
	var commandBuffer []byte
	for _, char := range command {
		commandBuffer = append(commandBuffer, byte(char))
	}

	var crlfBuffer [2]byte = [2]byte{'\r', '\n'}
	crlf := crlfBuffer[:]

	conn.Write(commandBuffer)
	conn.Write(crlf)
}
