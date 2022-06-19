package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"hostmonitor/arduino"
	"hostmonitor/grpc"
	"hostmonitor/measure"
	"hostmonitor/sensors"
	"hostmonitor/transport"
	"io/ioutil"
	"time"
)

const (
	BOARD_IP          = "127.0.0.1"
	PINGER_PROXY_PORT = 8090
)

type Configuration struct {
	VMSensor   bool   `yaml:"VMSensor"`
	CPUSensor  bool   `yaml:"CPU"`
	HostSensor bool   `yaml:"Host"`
	NetSensor  bool   `yaml:"NetSensor"`
	DiskSensor bool   `yaml:"Disk"`
	LoadSensor bool   `yaml:"Load"`
	Master     string `yaml:"Master"`
}

func loadConfiguration(path string) *Configuration {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not open %s error: %s\n", path, err)
		conf := &Configuration{true, true, true, true, true, true, "127.0.0.1:8758"}
		fmt.Printf("Host Monitor will use default configuration: %v\n", conf)
		return conf
	}
	if yfile == nil {
		panic("There was no error but YFile was null")
	}
	conf := Configuration{Master: "127.0.0.1:8758"}
	err2 := yaml.Unmarshal(yfile, &conf)
	if err2 != nil {
		fmt.Printf("Configuration file could not be parsed, error: %s\n", err2)
		panic(err2)
	}

	fmt.Printf("Found configuration: %v\n", conf)
	return &conf
}

func main() {
	conf := loadConfiguration("hostmonitor.yaml")

	reportCh := make(chan *measure.Measure)

	if conf.VMSensor {
		virtualMemorySensor := sensors.NewSensor(sensors.NewVirtualMemorySensor(time.Minute), "Disk Sensor", reportCh)
		virtualMemorySensor.Start()
	}
	if conf.CPUSensor {
		cpuSensor := sensors.NewSensor(sensors.NewCPUSensor(time.Minute), "CPU Sensor", reportCh)
		cpuSensor.Start()
	}
	if conf.HostSensor {
		hostSensor := sensors.NewSensor(sensors.NewHostSensor(time.Minute), "Host Sensor", reportCh)
		hostSensor.Start()
	}
	if conf.NetSensor {
		netSensor := sensors.NewSensor(sensors.NewNetSensor(time.Minute), "Net Sensor", reportCh)
		netSensor.Start()
	}
	if conf.DiskSensor {
		diskSensor := sensors.NewSensor(sensors.NewDiskSensor(time.Minute), "Disk Sensor", reportCh)
		diskSensor.Start()
	}
	// only works on linux
	if conf.LoadSensor {
		loadSensor := sensors.NewSensor(sensors.NewLoadSensor(time.Minute), "Load Sensor", reportCh)
		loadSensor.Start()
	}

	t := transport.NewUDPClient(conf.Master, reportCh, 60*time.Minute)
	t.Start()

	server := grpc.NewPingerProxy(PINGER_PROXY_PORT)
	server.Start()

	arduino := arduino.NewArduinoMonitor()
	arduino.Start(BOARD_IP)

	quitCh := make(chan int)
	<-quitCh
}
