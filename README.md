# host-monitor
## Description
This is the Host Monitor app, created by rFronteddu and maintained by cb218. It is a simple component designed to harvest host information such as Disk/VMemory/CPU status
The app is designed to run in conjunction with the **Testbed Monitor** (https://github.com/rFronteddu/testbed-monitor) app, another utility created by rFronteddu.

Sensors are imported from the library https://github.com/shirou/gopsutil.
Note that not all sensors produce results for every architecture. The Load Sensor, for example, only works in Linux.

## Configuration
The host monitor will look for a configuration file called **hostmonitor.yaml**.
Once a configuration file is found, active sensors must be defined. Sensor flags are used to enable each respective sensor.
If a conf file is not found, everything is active and statistics are delivered towards localhost.
Note that YAML is case-sensitive and there is no input validation.

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
## Connections
The host monitor listens for GRPC connections on port 8090.
Reports are sent to the Master address via UDP connection.

## Installation
### Prerequisites
* go > 1.6 
### Install
```
    go get 
    go build
```