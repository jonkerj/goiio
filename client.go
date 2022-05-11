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
	i.fetchContext()
	return i, nil
}

// bufferSizedReply sends a command and checks whether return response matches
// return length.
// The response expects 3 return arguments; response length, enabled channel
// buffers, and response
func (i *IIO) bufferSizedReply(cmd string) (*[]uint8, error) {
	log.Debugf("bufferSizedReply(%s)", cmd)
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

	var buffer []uint8

	if size == 0 {
		return &buffer, nil
	} else if size < 0 {
		return &buffer, fmt.Errorf("received negative size (error) from remote: %d", size)
	}

	buffer = make([]byte, size)

	// We do not care/check the channel buffers however we want to ensure it is
	// at least present to confirm that the entire response is in the expected
	// format
	_, err = i.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading channel size: %v", err)
	}

	numread, err := i.reader.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	if size != numread {
		return nil, fmt.Errorf("expected %d bytes, got %d", size, numread)
	}

	return &buffer, nil
}

// commandSizedReply sends a command and checks whether return response matches
// return length.
// The response expects two return arguments; response length and response.
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

	// If the response for commands such as OPEN, which returns only the size
	// ie. 0 with no response, this forever waits for the end of line character
	// so we therefore associate 0 as a successful response with no return
	// contents.
	reply := ""
	if size < 0 {
		return nil, fmt.Errorf("received negative size (error) from remote: %d", size)
	} else if size == 0 {
		return &reply, nil
	}

	reply, err = i.reader.ReadString('\n')
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

// CreateBuffer opens the specified device with the given mask of channels
func (i *IIO) CreateBuffer(devName string, bufSize int, chanMask string) error {

	_, err := i.commandSizedReply(fmt.Sprintf("OPEN %s %d %s", devName, bufSize, chanMask))
	if err != nil {
		return err
	}

	return nil
}

// DumpBuffer reads raw data from the specified device
func (i *IIO) DumpBuffer(devName string, bufSize int) ([]uint8, error) {

	var val *[]uint8
	val, err := i.bufferSizedReply(fmt.Sprintf("READBUF %s %d", devName, bufSize))
	if err != nil {
		return *val, err
	}

	return *val, nil
}

// DestroyBuffer closes the specified device
func (i *IIO) DestroyBuffer(devName string) error {

	_, err := i.commandSizedReply(fmt.Sprintf("CLOSE %s", devName))
	if err != nil {
		return err
	}

	return nil
}
