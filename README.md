# adbXchange

[![GoDoc](https://godoc.org/github.com/phlipse/adbxchange?status.svg)](https://godoc.org/github.com/phlipse/adbxchange)
[![Go Report Card](https://goreportcard.com/badge/github.com/phlipse/adbxchange)](https://goreportcard.com/report/github.com/phlipse/adbxchange)

adbXchange handles the exchange of ADB private keys. This is **useful if the user has to work with many devices with different keys**.

Also it transfers a workspace directory to /tmp/workspace/ on the target device and tries to make the binaries executable.

## Usage
Build adbXchange from source using the Makefile, rename the example configuration file to config.yml, adjust the configuration file content to your needs and run the executable.

If you want to change the default location of your configuration file (./config.yml) you can simply set the environment variable ADBXCHANGE_CONFIG accordingly.

**Error output**

If adbXchange does not start properly please have a look to the logfile. If you can not find anything in it please start adbXchange from console and have a look at the console output.

## License
[Apache License 2.0](https://github.com/phlipse/adbxchange/blob/master/LICENSE)
