# fxdmx [![CircleCI](https://circleci.com/gh/danesparza/fxdmx.svg?style=shield)](https://circleci.com/gh/danesparza/fxdmx)
REST service for DMX fixture control from Raspberry Pi.  

Control [lights](https://www.rollingstone.com/product-recommendations/lifestyle/best-stage-lights-928544/), [fog machines](https://www.amazon.com/dmx-fog-machine/s?k=dmx+fog+machine), [relays](https://www.amazon.com/ADJ-Products-Lighting-Dimmer-DP-415R/dp/B07C7Y4MT9/ref=sr_1_10?dchild=1&keywords=dmx+relay&qid=1626704824&sr=8-10) ... even [flame throwers](https://www.youtube.com/watch?v=jbIG1ijw9Qw)! 

 Made with ❤️ for makers, DIY craftsmen, prop makers and professional soundstage designers everywhere
 
 ## Prerequisites
There are no other software prerequisites, but you'll need to make sure you have a USB DMX controller.  I recommend the [DMXking ultraDMX micro](https://dmxking.com/usbdmx/ultradmxmicro) -- also available at Amazon.  You should probably also pick up a [few lights](https://www.amazon.com/gp/product/B07DPGPRZ3/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1) and at least 1 [3 prong DMX cable](https://www.amazon.com/gp/product/B0885HHY5Q/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1).  
## Installing
Installing fxdmx is also really simple.  Grab the .deb file from the [latest release](https://github.com/danesparza/fxdmx/releases/latest) and then install it using dpkg:


```bash
sudo dpkg -i fxdmx-1.0.40_armhf.deb 
````

This automatically installs the **fxdmx** service with a default configuration and starts the service. 

You can then use the service at http://localhost:3040

See the REST API documentation at http://localhost:3040/v1/swagger/

## Setup
There is one very simple setup step you should probably do:  Setting the default serial USB device.  You can see all the USB serial devices installed by using the REST service call `/v1/system/usbinfo`.  On my test Raspberry Pi, here's what this looks like when I run curl:

Request:
```bash
curl -X GET "http://localhost:3040/v1/system/usbinfo" -H  "accept: application/json"
```
Response:
```json
{
  "message": "1 devices found",
  "data": [
    {
      "device": "/dev/ttyUSB0",
      "product": "DMX USB PRO",
      "manufacturer": "DMXking.com"
    }
  ]
}
```
Notice that 'device' property that shows a path like `/dev/ttyUSB0`?  You'll need to take what you find there and navigate to the REST service call `/v1/system/defaultusb` to set the default device to use.  Here's what it looks like for me using curl:

Request:
```bash
curl -X PUT "http://localhost:3040/v1/system/defaultusb" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"devicepath\": \"/dev/ttyUSB0\"}"
```
Response:
```json
{
  "message": "Default USB device updated",
  "data": "/dev/ttyUSB0"
}
```
Now you can run your DMX timelines without having to set the device information every time.

## Removing 
Uninstalling is just as simple:

```bash
sudo dpkg -r fxdmx
```
