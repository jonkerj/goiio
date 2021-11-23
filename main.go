package main

import (
	"github.com/jonkerj/go-iio/pkg/client"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	c, err := client.New("neon.thuis.efgh.nl:30431")
	if err != nil {
		panic(err)
	}

	log.Debugf("%d devices found", len(c.Context.Devices))

	for _, dev := range c.Context.Devices {
		log.Infof("Device: id=%s, name=%s\n", dev.ID, dev.Name)
		for _, ch := range dev.Channels {
			log.Infof("  Channel: id=%s\n", ch.ID)
			for _, attr := range ch.Attributes {
				log.Infof("    Attribute: %s\n", attr.Name)
			}
		}
	}
}
