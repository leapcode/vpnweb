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
	"0xacab.org/leap/vpnweb/pkg/auth/creds"
	"fmt"
	"github.com/reiver/go-telnet"
	"log"
	"time"
)

const (
	Label                 string = "sip2"
	loginRequestTemplate  string = "9300CN%s|CO%s|CP%s|"
	statusRequestTemplate string = "23000%s    %sAO%s|AA%s|AD%s|"
)

type sipClient struct {
	host     string
	port     string
	location string
	conn     *telnet.Conn
	parser   *Parser
}

func newClient(host, port, location string) sipClient {
	c := sipClient{host, port, location, nil, nil}
	c.parser = getParser()
	return c
}

func (c *sipClient) Connect() (bool, error) {
	conn, err := telnet.DialTo(c.host + ":" + c.port)
	if nil != err {
		log.Println("error", err)
		return false, err
	}
	c.conn = conn
	return true, nil
}

func (c *sipClient) Login(user, pass string) bool {
	loginStr := fmt.Sprintf(loginRequestTemplate, user, pass, c.location)
	if nil == c.conn {
		fmt.Println("error! null connection")
	}
	telnetSend(c.conn, loginStr)
	loginResp := telnetRead(c.conn)
	msg := c.parseResponse(loginResp)
	if value, ok := c.parser.getFixedFieldValue(msg, okVal); ok && value == trueVal {
		return true
	}
	return false
}

func (c *sipClient) parseResponse(txt string) *message {
	msg := c.parser.parseMessage(txt)
	return msg
}

/* Authenticator interface */

func (c *sipClient) GetLabel() string {
	return Label
}

func (c *sipClient) NeedsCredentials() bool {
	return true
}

func (c *sipClient) CheckCredentials(credentials *creds.Credentials) bool {
	currentTime := time.Now()
	user := credentials.User
	passwd := credentials.Password
	statusRequest := fmt.Sprintf(
		statusRequestTemplate,
		currentTime.Format("20060102"),
		currentTime.Format("150102"),
		c.location, user, passwd)
	telnetSend(c.conn, statusRequest)

	statusMsg := c.parseResponse(telnetRead(c.conn))
	if value, ok := c.parser.getFieldValue(statusMsg, validPatron); ok && value == yes {
		if value, ok := c.parser.getFieldValue(statusMsg, validPatronPassword); ok && value == yes {
			return true
		}
	}

	// TODO log whatever error we can find (AF, Screen Message, for instance)
	return false
}
