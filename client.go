package goiio

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Creates a new IIO client
func New(remote string) (*IIO, error) {
	log.Debugf("resolving tcp to %s", remote)
	addr, err := net.ResolveTCPAddr("tcp", remote)
	if err != nil {
		return nil, fmt.Errorf("error resolving TCP address: %v", err)
	}

	log.Debugf("dialing %v", addr)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to remote: %v", err)
	}

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	i := &IIO{conn: conn, reader: r, writer: w}

	err = i.fetchContext()
	if err != nil {
		return nil, fmt.Errorf("error fetching context: %v", err)
	}

	return i, nil
}

func (i *IIO) commandSizedReply(cmd string) (*string, error) {
	log.Debugf("commandSizedReply(%s)", cmd)
	_, err := i.writer.WriteString(fmt.Sprintf("%s\n", cmd))
	if err != nil {
		return nil, fmt.Errorf("error writing: %v", err)
	}

	if err := i.writer.Flush(); err != nil {
		return nil, fmt.Errorf("error flushing: %v", err)
	}

	log.Debugf("reading reply")
	sizeStr, err := i.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading size: %v", err)
	}
	sizeStr = strings.TrimRight(sizeStr, "\n")

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, fmt.Errorf("error converting size '%s' into int: %v", sizeStr, err)
	}
	log.Debugf("size is %d", size)

	if size < 0 {
		return nil, fmt.Errorf("received negative size (error) from remote: %d", size)
	}

	reply, err := i.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}
	reply = strings.TrimRight(reply, "\n")
	log.Debugf("reply is %s", reply)

	if size != len(reply) {
		return nil, fmt.Errorf("expected %d bytes, got %d", size, len(reply))
	}

	return &reply, nil
}

func (i *IIO) fetchContext() error {
	var ctx Context

	contextStr, err := i.commandSizedReply("PRINT")
	if err != nil {
		return err
	}

	if err := xml.Unmarshal([]byte(*contextStr), &ctx); err != nil {
		return fmt.Errorf("error unmarshalling XML: %v", err)
	}

	i.Context = &ctx

	return nil
}

func (i *IIO) FetchAttributes() error {
	for _, device := range i.Context.Devices {
		for _, channel := range device.Channels {
			if channel.Type != "input" {
				continue
			}

			for _, attribute := range channel.Attributes {
				raw, err := i.commandSizedReply(fmt.Sprintf("READ %s INPUT %s %s", device.ID, channel.ID, attribute.Name))
				if err != nil {
					return fmt.Errorf("error reading input: %v", err)
				}

				rawCleaned := strings.TrimRight(*raw, "\x00")

				f, err := strconv.ParseFloat(rawCleaned, 64)
				if err != nil {
					return fmt.Errorf("reply is not parseble as float: %v", err)
				}

				attribute.Value = f
			}
		}
	}

	return nil
}
