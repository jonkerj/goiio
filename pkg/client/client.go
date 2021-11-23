package client

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/jonkerj/go-iio/pkg/types"
	log "github.com/sirupsen/logrus"
)

type IIO struct {
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	Context *types.Context
}

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
	i.fetchContext()
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
	var ctx types.Context

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
