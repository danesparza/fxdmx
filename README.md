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


## Network discovery (fxcontroller)

The service advertises its HTTP API through Zeroconf (mDNS/DNS-SD) by
default. A future fxcontroller can browse **`_fx._tcp` in `local.`** once
to find fxaudio, fxpixel, fxdmx, and fxtrigger. No controller is required
to run the services, and direct HTTP API access remains available if
multicast registration fails (a warning is logged).

```yaml
discovery:
  enabled: true
  name: "" # Optional friendly instance name, 1–63 UTF-8 bytes
  id: ""   # Optional unique, stable installation ID
```

The default name is `service-hostname-port`; the default ID is
`service:hostname:port`. Give each installation a distinct `discovery.id`
if identity must survive hostname or port changes. Names must also be
unique on the LAN. TXT values must fit the DNS-SD 255-byte record limit.

The discovery contract is the same across all four projects:

| DNS record | Meaning |
| --- | --- |
| SRV / A / AAAA | HTTP hostname, actual listening port, and addresses |
| TXT `txtvers=1` | Discovery metadata format version |
| TXT `service` | `fxaudio`, `fxpixel`, `fxdmx`, or `fxtrigger` |
| TXT `id` | Installation identity |
| TXT `api=v1` | API version |
| TXT `scheme=http` | API transport |
| TXT `path=/v1` | API base path |

Controllers should use the resolved SRV port and address, dispatch by
`service`, check supported metadata/API versions, and verify API reachability
separately. Discovery is not authentication; treat advertisements as untrusted
network input. This adds service advertisements, not a controller or UI.

Advertisements start only after the HTTP listener binds and are withdrawn
on SIGINT/SIGTERM or when the HTTP server stops. Discovery uses UDP 5353
multicast on the local network segment; firewalls, Wi-Fi client isolation,
VLANs, and container networking can prevent discovery. Containers typically
need host networking (where supported) or an mDNS relay. No internet service
or separately installed Bonjour/Avahi daemon is required.

On macOS, inspect live advertisements with:

```sh
dns-sd -B _fx._tcp local.
# Resolve one of the instance names returned above:
dns-sd -L "INSTANCE NAME" _fx._tcp local.
```

On Linux with Avahi tools: `avahi-browse -rt _fx._tcp`.
