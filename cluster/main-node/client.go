package client

import (
	"bufio"
	"net"
	"strings"
	"time"
) 

type Client struct {
	conn  net.Conn 
	reader *bufio.Reader
} 

func NewClient(address string) (*Client , error) {
	conn, err := net.DialTimeout("tcp" , address , 2*time.Second)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)

	// Read the entire banner and prompt
	var welcome strings.Builder
	for {
		chunk, err := reader.ReadString('>')
		if err != nil {
			return nil, err
		}
		welcome.WriteString(chunk)

		if strings.Contains(welcome.String(), "----------") {
			break
		}
	}

	return &Client{
		conn:   conn,
		reader: reader,
	}, nil
}

func (c *Client) SendCommand(cmd string) (string ,error) {
	_, err := c.conn.Write([]byte(cmd +"\n")) 
	if err != nil{
		return "" , err
 	} 

	var response strings.Builder 
	for {
		line , err := c.reader.ReadString('>') 
		if err != nil {
			return "" , err 
		} 
		response.WriteString(line) 
		if strings.HasSuffix(line, "real-db>") {
			break
		}
	} 
	return strings.TrimSuffix(response.String() ,"real-db>"), nil 
} 

func (c *Client) Close(){
	c.conn.Close()
}

