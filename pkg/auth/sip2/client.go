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
	"errors"
	"fmt"
	"github.com/linxiaozhi/go-telnet"
	"log"
	"time"
)

const (
	Label                 string = "sip2"
	loginRequestTemplate  string = "9300CN%s|CO%s|CP%s|"
	statusRequestTemplate string = "23000%s    %sAO%s|AA%s|AD%s|"
	scStatusRequest       string = "9901002.00"
	heartBeatSeconds             = 240
)

type sipClient struct {
	host          string
	port          string
	location      string
	user          string
	pass          string
	conn          gote.Connection
	heartBeatDone chan bool
	reqQueue      chan request
	parser        *Parser
}

type request struct {
	msg      string
	respChan chan response
}

type response struct {
	msg string
	err error
}

func newClient(host, port, location string) sipClient {
	reqQ := make(chan request)
	parser := getParser()
	c := sipClient{host, port, location, "", "", nil, nil, reqQ, parser}
	return c
}

func (c *sipClient) startDispatcher() {
	go func() {
		for {
			req := <-c.reqQueue
			resp, err := c.handleRequest(req.msg)
			req.respChan <- response{resp, err}
		}
	}()
}

func (c *sipClient) startHeartBeat() {
	ticker := time.NewTicker(heartBeatSeconds * time.Second)
	c.heartBeatDone = make(chan bool)
	go func() {
		for {
			select {
			case <-c.heartBeatDone:
				log.Println("Stopping heartbeat")
				return
			case <-ticker.C:
				resp, err := c.sendRequest(scStatusRequest)
				/* TODO
				   for now we are only interested in letting the server
				   know that we are alive.
				   but we could parse the response fields */
				if err != nil {
					log.Println(">> status error:", err)
				}
				log.Println(">> sip status:", resp)
			}
		}
	}()
}

func (c *sipClient) sendRequest(msg string) (string, error) {
	respChan := make(chan response)
	c.reqQueue <- request{msg, respChan}
	resp := <-respChan
	return resp.msg, resp.err
}

func (c *sipClient) handleRequest(msg string) (string, error) {
	err := telnetSend(c.conn, msg)
	if err != nil {
		return "", err
	}
	resp, err := telnetRead(c.conn)
	return resp, err
}

func (c *sipClient) doConnectAndReact() (bool, error) {
	_, err := c.connect()
	if err != nil {
		return false, err
	}
	_, err = c.login()
	if err != nil {
		return false, err
	}
	c.startDispatcher()
	c.startHeartBeat()
	return true, nil

}

func (c *sipClient) setCredentials(user, pass string) {
	c.user = user
	c.pass = pass
}

func (c *sipClient) connect() (bool, error) {
	conn, err := gote.DialTimeout("tcp", c.host+":"+c.port, time.Second*2)
	if nil != err {
		log.Println("connect error:", err)
		return false, err
	}
	c.conn = conn
	return true, nil
}

func (c *sipClient) login() (bool, error) {
	loginStr := fmt.Sprintf(loginRequestTemplate, c.user, c.pass, c.location)
	if nil == c.conn {
		fmt.Println("error! null connection")
		return false, errors.New("null connection received")
	}

	err := telnetSend(c.conn, loginStr)
	if err != nil {
		log.Println("Error while sending request")
		return false, err
	}

	loginResp, err := telnetRead(c.conn)
	if err != nil {
		log.Println("login error on read: ", err)
		return false, err
	}

	msg, err := c.parseResponse(loginResp)
	if err != nil {
		log.Println("login error:", err)
		return false, err
	}
	if value, ok := c.parser.getFixedFieldValue(msg, okVal); ok && value == trueVal {
		log.Println("SIP admin login OK")
		return true, nil
	}
	return false, nil
}

func (c *sipClient) close() {
	c.heartBeatDone <- true
	c.conn.Close()
}

func (c *sipClient) parseResponse(txt string) (*message, error) {
	msg, err := c.parser.parseMessage(txt)
	return msg, err
}

/* Authenticator interface */

func (c *sipClient) GetLabel() string {
	return Label
}

func (c *sipClient) NeedsCredentials() bool {
	return true
}

func (c *sipClient) CheckCredentials(credentials *creds.Credentials) (bool, error) {

	currentTime := time.Now()
	user := credentials.User
	passwd := credentials.Password

	statusRequest := fmt.Sprintf(
		statusRequestTemplate,
		currentTime.Format("20060102"),
		currentTime.Format("150102"),
		c.location, user, passwd)

	statusMsg := &message{}
	resp, err := c.sendRequest(statusRequest)
	if err != nil {
		return false, err
	}

	statusMsg, err = c.parseResponse(resp)
	if err != nil {
		log.Println("Error while parsing response")
		return false, err
	}
	if value, ok := c.parser.getFieldValue(statusMsg, validPatron); ok && value == yes {
		if value, ok := c.parser.getFieldValue(statusMsg, validPatronPassword); ok && value == yes {
			return true, nil
		}
		// TODO log whatever error we can find (AF, Screen Message, for instance)
	}

	return false, errors.New("unknown error while checking credentials")
}
