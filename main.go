package main

import "github.com/danesparza/fxdmx/cmd"

// @title fxDmx
// @version 1.0
// @description fxDmx REST service for DMX fixture control from Raspberry Pi

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
func main() {
	cmd.Execute()
}
