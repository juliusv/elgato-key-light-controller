# Elgato Key Light Controller

This is a tiny Go program to control the power, brightness, and temperature settings for Elgato Key Lights with known IP addresses. It was hacked together in ~1h and works for my purposes, but feel free to reuse and/or improve.

## Building

```bash
go build
```

## Running

```bash
./elgato-key-light-controller [...]
```

For example:

```bash
./elgato-key-light-controller -light-ips="192.168.0.181,192.168.0.182" -command=toggle-power
```

## Flags

```bash
$ ./elgato-key-light-controller -help
Usage of ./elgato-key-light-controller:
  -command string
    	Command to run. May be: toggle-power, decrease-brightness, increase-brightness, decrease-temperature, increase-temperature, set-min-brightness, set-max-brightness, set-min-temperature, set-max-temperature, set-brightness, set-temperature. (default "toggle-power")
  -light-ips string
    	Comma-separated list of Elgato Key Light IPs. (default "192.168.0.181,192.168.0.182")
  -value string
    	Numeric value to use for 'set-brightness' and 'set-temperature' commands.
```
