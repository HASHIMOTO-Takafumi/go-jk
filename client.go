package jk

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
)

type Client struct {
	address    string
	connection net.Conn
	tp         *textproto.Reader
}

func NewClient(address string, option string) (*Client, error) {
	var err error
	self := new(Client)
	self.address = address

	self.connection, err = net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	//     defer self.connection.Close()

	r := bufio.NewReader(self.connection)
	self.tp = textproto.NewReader(r)

	fmt.Fprintf(self.connection, option)
	for {
		line, _ := self.tp.ReadLine()
		if strings.Index(line, "OK") != -1 {
			break
		}
	}
	return self, err
}

func (self *Client) RawParse(query string) ([]string, error) {
	fmt.Fprintf(self.connection, "%s\n", query)

	lines := make([]string, 0)

	for {
		line, err := self.tp.ReadLine()
		if err != nil {
			return nil, err
		} else if line == "EOS" {
			break
		}
		lines = append(lines, line)
	}
	return lines, nil
}
