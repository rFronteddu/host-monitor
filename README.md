# host-monitor
## Description
This is the Host Monitor app, created by rFronteddu and maintained by cb218. It is a simple component designed to harvest host information such as Disk/VMemory/CPU status.

The app is designed to run in conjunction with the **Testbed Monitor** (https://github.com/rFronteddu/testbed-monitor) app, another utility created by rFronteddu.

Sensors are imported from the library https://github.com/shirou/gopsutil.
Note that not all sensors produce results for every architecture. The Load Sensor, for example, only works in Linux.

## Inputs
The host monitor will look for a configuration file called **hostmonitor.yaml**.If no conf file is  found, everything is active and statistics are delivered towards localhost.

### hostmonitor.yaml
```
    # Enables the virtual memory sensor (RAM reports)
    VMSensor: true
    # Enables the CPU usage (% CPU used)
    CPU: true
    # Enables the HOST report (OS, CPU, etc)
    Host: true
    # Enables the NETWORK report (Bytes generated and received by the host)
    NetSensor: false
    # Enables the DISK report (HD free space)
    Disk: true
    # Enables the Load reports (1/5/15 reports, only works on unix)
    Load: false
    # Specifies where to send reports, IP address and Port
    Master: 127.0.0.1:8758
    # The IP of the controller to ping periodically for reachability
    # (if not defined, it will be obtained as the host IP x.x.x.1)
    #BoardIP: ""
    # How often a report is sent in minutes (default 30)
    ReportPeriod: "2"
    # The port to listen for GRPC pings
    PingerProxyPort: "8100"
    # MQTT subscription information'
    MQTTBroker: ""
    MQTTTopic: ""
```

## Installation
### Prerequisites
* go > 1.6 
### Install
If using a linux system, the **host-monitor.service** daemon can be used to build and run the program.
```
    ./setup-service.sh
```
On other systems, the program can be started manually after building.
```
    go build main.go
    ./main
```
