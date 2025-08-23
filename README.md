# fxdmx [![Build and release](https://github.com/danesparza/fxdmx/actions/workflows/release.yaml/badge.svg)](https://github.com/danesparza/fxdmx/actions/workflows/release.yaml)
REST service for DMX fixture control from Raspberry Pi.  

Control [lights](https://www.rollingstone.com/product-recommendations/lifestyle/best-stage-lights-928544/), [fog machines](https://www.amazon.com/dmx-fog-machine/s?k=dmx+fog+machine), [relays](https://www.amazon.com/ADJ-Products-Lighting-Dimmer-DP-415R/dp/B07C7Y4MT9/ref=sr_1_10?dchild=1&keywords=dmx+relay&qid=1626704824&sr=8-10) ... even [flame throwers](https://www.youtube.com/watch?v=jbIG1ijw9Qw)! 

 Made with ❤️ for makers, DIY craftsmen, prop makers and professional soundstage designers everywhere

## Installation
### Prerequisites
Install Raspberry Pi OS (Bookworm or later) on your device. For best results, use the [Raspberry Pi Imager](https://www.raspberrypi.com/software/) and select **Raspberry Pi OS (64-bit)**.

Install the prerequisite package repository (one-time setup per machine):

``` bash
wget https://packages.cagedtornado.com/prereq.sh -O - | sh
```

### Installing fxdmx
Install the fxdmx package:

``` bash
sudo apt install fxdmx
```

You can then use the service at http://localhost:3040

See the REST API documentation at http://localhost:3040/v1/swagger/

## Setup
After plugging in your hardware, there is one very simple setup step you should probably do:  Setting the default serial USB device.  You can see all the USB serial devices installed by using the REST service call `/v1/system/usbinfo`.  On my test Raspberry Pi, here's what this looks like when I run curl:

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
