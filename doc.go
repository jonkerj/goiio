// goiio is a pure golang client for IIO (industrial I/O)
//
// Introduction
//
// go-iio was written to access remote iiod instanced, since libiio does
// not seem to have golang bindings by itself.
// It's main use case is connecting with remote IIOd instances: local IIO
// access is out of scope, you'll need to run iiod.
//
// See https://github.com/analogdevicesinc/libiio
//
// Example
//
// The following example shows an example on how to use this library
//
//   c, err := goiio.New("my-sensor.home.lan:30431")
//   if err != nil {
//     panic(err)
//   }
//
//   // Populate the values in attributes
//   if err = c.FetchAttributes(); err != nil {
//     panic(err)
//   }
//
//   //
//   for _, dev := range c.Context.Devices {
//     log.Infof("Device: id=%s, name=%s", dev.ID, dev.Name)
//     for _, ch := range dev.Channels {
//       log.Infof("  Channel: id=%s", ch.ID)
//       for _, attr := range ch.Attributes {
//         log.Infof("    Attribute: %s, value: %0.3f", attr.Name, attr.Value)
//       }
//     }
//   }
package goiio
