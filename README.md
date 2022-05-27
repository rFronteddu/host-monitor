# host-monitor
The host monitor is a simple component designed to harvest host information such as Disk/VMemory/CPU status. Active sensors can be configured by writing a yml file. Run hostmonitor -h to see options.

## GRPC
The host monitor listen for GRPC connections on port 8090.

## Configuration
The host monitor will look for a configuration file in the following locations:
* hostmonitor.yaml
* /usr/lib/sensei/yaml/hostmonitor/hostmonitor.yaml

Once a configuration file is found, active sensors must be defined.

### Configuration example
```
    # Enables the virtual memory sensor (RAM reports)
    VMSensor: true
    # Enables the CPU usage (% CPU used)
    CPU: true
    # Enables the HOST report (OS, CPU, etc)
    Host: true
    # Enables the NETWORK report (Bytes generated and received by the host)
    NetSensor: true
    # Enables the DISK report (HD free space)
    Disk: true 
    # Enables the Load reports (1/5/15 reports, only works on unix)
    Load: true
    # Specifies where to send reports
    Master: 127.0.0.1:8758
```

## Installation
### Prerequisites
* go > 1.6 
### Install
```
    go get 
    go build.
```