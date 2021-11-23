package types

import "encoding/xml"

type Context struct {
	XMLName          xml.Name           `xml:"context"`
	Text             string             `xml:",chardata"`
	Name             string             `xml:"name,attr"`
	Description      string             `xml:"description,attr"`
	ContextAttribute []ContextAttribute `xml:"context-attribute"`
	Devices          []Device           `xml:"device"`
}

type ContextAttribute struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Device struct {
	Text     string    `xml:",chardata"`
	ID       string    `xml:"id,attr"`
	Name     string    `xml:"name,attr"`
	Channels []Channel `xml:"channel"`
}

type Channel struct {
	Text       string             `xml:",chardata"`
	ID         string             `xml:"id,attr"`
	Type       string             `xml:"type,attr"`
	Attributes []ChannelAttribute `xml:"attribute"`
}

type ChannelAttribute struct {
	Text     string `xml:",chardata"`
	Name     string `xml:"name,attr"`
	Filename string `xml:"filename,attr"`
}
