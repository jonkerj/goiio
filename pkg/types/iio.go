package types

import "encoding/xml"

// Context holds the context of the remote IIOd
type Context struct {
	XMLName          xml.Name           `xml:"context"`           // Not sure where this is for
	Text             string             `xml:",chardata"`         // Probably not used
	Name             string             `xml:"name,attr"`         // The name of the context
	Description      string             `xml:"description,attr"`  // Description of the context
	ContextAttribute []ContextAttribute `xml:"context-attribute"` // Attributes that apply to the context
	Devices          []Device           `xml:"device"`            // Devices connected. This is probably what you need
}

// ContextAttribute is the data structure containing attributes relating to the context
type ContextAttribute struct {
	Text  string `xml:",chardata"`  // Probably not used, this is how XML works
	Name  string `xml:"name,attr"`  // Name of the attribute, e.g. "local,kernel"
	Value string `xml:"value,attr"` // Value of the attribute, e.g. "5.11.0-40-generic"
}

// Device holds information regarding IIO device
type Device struct {
	Text     string    `xml:",chardata"` // Probably not used, this is how XML works
	ID       string    `xml:"id,attr"`   // Identifier, e.g. "iio:device0"
	Name     string    `xml:"name,attr"` // Name of the device, e.g. "bme280"
	Channels []Channel `xml:"channel"`   // List of channels in the device
}

// Channels holds all information related to a channel in a IIO device
type Channel struct {
	Text       string             `xml:",chardata"` // Probably not used, this is how XML works
	ID         string             `xml:"id,attr"`   // ID of the channel, e.g. "temp"
	Type       string             `xml:"type,attr"` // Type of the channel, e.g. "input"
	Attributes []ChannelAttribute `xml:"attribute"` // List of attributes of the channel
}

// ChannelAttribute is the data structure describing an attribute of a channel
type ChannelAttribute struct {
	Text     string `xml:",chardata"`     // Probably not used, this is how XML works
	Name     string `xml:"name,attr"`     // Name of the attribute, e.g. "oversampling_ratio"
	Filename string `xml:"filename,attr"` // Filename of the attribute in sysfs, e.g. "in_pressure_oversampling_ratio"
}
