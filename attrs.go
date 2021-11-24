package goiio

// CompareDevice is a type passed to Context.GetDevice
// implementing some kind of comparison function
type CompareDevice func(*Device) bool

// CompareContextAttribute is a type passed to Context.GetAttribute
// implementing some kind of comparison function
type CompareContextAttribute func(*ContextAttribute) bool

// CompareChannel is a type passed to Device.GetChannel
// implementing some kind of comparison function
type CompareChannel func(*Channel) bool

// CompareChannelAttribute is a type passed to Channel.GetAttribute
// implementing some kind of comparison function
type CompareChannelAttribute func(*ChannelAttribute) bool

// GetDevice fetches the first device from the context satisfying the comparison
// function
func (h *Context) GetDevice(comp CompareDevice) *Device {
	for _, device := range h.Devices {
		if comp(device) {
			return device
		}
	}

	return nil
}

// GetDeviceByName fetches the first device with given name
func (c *Context) GetDeviceByName(name string) *Device {
	return c.GetDevice(func(d *Device) bool { return d.Name == name })
}

// GetDeviceByID fetches the first device with given name
func (c *Context) GetDeviceByID(id string) *Device {
	return c.GetDevice(func(d *Device) bool { return d.ID == id })
}

// GetDeviceByNameOrID fetches the first device that does not mismatch given
// name or id. If only name is given, it behaves like GetDeviceByName and when
// only id is given like GetDeviceByID. When both are given, it matches on
// both attributes and when none are given, it matches any first device.
func (c *Context) GetDeviceByNameOrID(name, id *string) *Device {
	return c.GetDevice(func(d *Device) bool {
		if name != nil && d.Name != *name {
			return false
		}
		if id != nil && d.ID != *id {
			return false
		}
		return true
	})
}

// GetAttribute fetches the first attribute from the context satisfying the comparison
// function
func (c *Context) GetAttribute(comp CompareContextAttribute) *ContextAttribute {
	for _, attr := range c.ContextAttributes {
		if comp(attr) {
			return attr
		}
	}

	return nil
}

// GetAttributeByName fetches the first attribute with given name
func (c *Context) GetAttributeByName(name string) *ContextAttribute {
	return c.GetAttribute(func(c *ContextAttribute) bool { return c.Name == name })
}

// GetChannel fetches the first channel from the device satisfying the comparison
// function
func (d *Device) GetChannel(comp CompareChannel) *Channel {
	for _, ch := range d.Channels {
		if comp(ch) {
			return ch
		}
	}

	return nil
}

// GetChannelByID fetches the first attribute with given ID
func (d *Device) GetChannelByID(id string) *Channel {
	return d.GetChannel(func(c *Channel) bool { return c.ID == id })
}

// GetAttribute fetches the first attribute from the channel satisfying the comparison
// function
func (c *Channel) GetAttribute(comp CompareChannelAttribute) *ChannelAttribute {
	for _, attr := range c.Attributes {
		if comp(attr) {
			return attr
		}
	}

	return nil
}

// GetAttrByName fetches the first attribute with given name
func (c *Channel) GetAttrByName(name string) *ChannelAttribute {
	return c.GetAttribute(func(ca *ChannelAttribute) bool { return ca.Name == name })
}
