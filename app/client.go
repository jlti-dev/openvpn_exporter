package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	Host      string
	Port      string
	cmdActive sync.Mutex
	conn      net.Conn
	reader    *bufio.Reader
	offline   bool
}

func NewClient(host string, port string) (*Client, error) {
	c := &Client{}
	c.Host = host
	c.Port = port

	return c, c.connect()

}

func (c *Client) connect() error {
	connUrl := fmt.Sprintf("%s:%s", c.Host, c.Port)

	log.Printf("Connecting to %s", connUrl)
	conn, err := net.Dial("tcp", connUrl)
	if err != nil {
		return err
	}
	log.Printf("%s: Reading first line (welcome)", c.Host)
	c.reader = bufio.NewReader(conn)
	line, errRead := c.reader.ReadString('\n') //read welcome message
	if errRead != nil {
		return errRead
	}
	log.Printf("%s: %s", c.Host, line)
	c.conn = conn
	return nil
}
func (c *Client) GetVersion() (*Version, error) {
	ans, err := c.execute("version")
	if err != nil {
		return nil, err
	}
	return parseVersion(ans)
}
func (c *Client) GetStats() (*LoadStats, error) {
	ans, err := c.execute("load-stats")
	if err != nil {
		return nil, err
	}
	return parseStats(ans)
}
func (c *Client) GetDetails() (*Status, error) {
	ans, err := c.execute("status 2")
	if err != nil {
		return nil, err
	}
	return parseStatus(ans)
}

//Execute connects to the OpenVPN server, sends command and reads response
func (c *Client) execute(cmd string) (string, error) {
	log.Println(c.conn)
	//log.Printf("%s: Sending command: %s\n", c.Host, cmd)
	if _, err := fmt.Fprintf(c.conn, cmd+"\n"); err != nil {
		log.Printf("%s: Error: %s\n", err, c.Host)
		return "", c.connect()
	}
	log.Printf("%s: Command sent successful: %s\n", c.Host, cmd)

	return c.readResponse()
}

//ReadResponse
func (c *Client) readResponse() (string, error) {
	var finished = false
	var result = ""

	for !finished {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			log.Printf("%s: %s (%s)", c.Host, line, err)
			return "", err
		}

		result += line
		if strings.Index(line, "END") == 0 ||
			strings.Index(line, "SUCCESS:") == 0 ||
			strings.Index(line, "ERROR:") == 0 {
			finished = true
		}
	}
	if os.Getenv("DEBUG") == "true" {
		log.Printf("%s: %s", c.Host, result)
	}
	return result, nil
}
