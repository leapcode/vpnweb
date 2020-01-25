package sip2

import (
	"fmt"
	"github.com/reiver/go-telnet"
	"log"
	"time"
)

const loginRequestTemplate string = "9300CN%s|CO%s|CP%s|"
const statusRequestTemplate string = "23000%s    %sAO%s|AA%s|AD%s|"

type Client struct {
	Host     string
	Port     string
	location string
	conn     *telnet.Conn
	parser   *Parser
}

func NewClient(host, port, location string) Client {
	c := Client{host, port, location, nil, nil}
	c.parser = getParser()
	return c
}

func (c *Client) Connect() (bool, error) {
	conn, err := telnet.DialTo(c.Host + ":" + c.Port)
	if nil != err {
		log.Println(log.Printf("error: %v", err))
		return false, err
	}
	c.conn = conn
	return true, nil
}

func (c *Client) Login(user, pass string) bool {
	loginStr := fmt.Sprintf(loginRequestTemplate, user, pass, c.location)
	if nil == c.conn {
		fmt.Println("error! null connection")
	}
	telnetSend(c.conn, loginStr)
	loginResp := telnetRead(c.conn)
	msg := c.parseResponse(loginResp)
	if value, ok := c.parser.getFixedFieldValue(msg, Ok); ok && value == TRUE {
		return true
	}
	return false
}

func (c *Client) CheckCredentials(user, passwd string) bool {
	currentTime := time.Now()
	statusRequest := fmt.Sprintf(
		statusRequestTemplate,
		currentTime.Format("20060102"),
		currentTime.Format("150102"),
		c.location, user, passwd)
	telnetSend(c.conn, statusRequest)

	statusMsg := c.parseResponse(telnetRead(c.conn))
	if value, ok := c.parser.getFieldValue(statusMsg, ValidPatron); ok && value == YES {
		if value, ok := c.parser.getFieldValue(statusMsg, ValidPatronPassword); ok && value == YES {
			return true
		}
	}

	// TODO log whatever error we can find (AF, Screen Message, for instance)
	return false
}

func (c *Client) parseResponse(txt string) *Message {
	msg := c.parser.parseMessage(txt)
	return msg
}
