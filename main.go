package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"hostmonitor/grpc"
	"hostmonitor/measure"
	"hostmonitor/mqtt"
	"hostmonitor/probers"
	"hostmonitor/sensors"
	"hostmonitor/transport"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
)

type Configuration struct {
	VMSensor        bool   `yaml:"VMSensor"`
	CPUSensor       bool   `yaml:"CPU"`
	HostSensor      bool   `yaml:"Host"`
	NetSensor       bool   `yaml:"NetSensor"`
	DiskSensor      bool   `yaml:"Disk"`
	LoadSensor      bool   `yaml:"Load"`
	Master          string `yaml:"Master"`
	BoardIP         string `yaml:"BoardIP"`
	ReportPeriod    string `yaml:"ReportPeriod"`
	PingerProxyPort string `yaml:"PingerProxyPort"`
	MQTTBroker      string `yaml:"MQTTBroker"`
	MQTTTopic       string `yaml:"MQTTTopic"`
}

func loadConfiguration(path string) *Configuration {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not open %s error: %s\n", path, err)
		conf := &Configuration{true, true, true, true, true, true, "127.0.0.1:8758", "127.0.0.1", "30", "8090", "", ""}
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

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	for {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			fmt.Println(err.Error() + ", will try again in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()

		return conn.LocalAddr().(*net.UDPAddr).IP.String()
	}
}

func main() {
	conf := loadConfiguration("hostmonitor.yaml")
	reportCh := make(chan *measure.Measure)

	var boardAddress string
	if conf.BoardIP != "" {
		boardAddress = conf.BoardIP
	} else {
		outboundIP := GetOutboundIP()
		s := strings.Split(outboundIP, ".")
		s[len(s)-1] = "1"
		boardAddress = strings.Join(s[:], ".")
		fmt.Println("The board address was not specified so it was automatically detected as: " + boardAddress)
	}
	board := probers.NewBoardMonitor()
	board.Start(boardAddress, reportCh)

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

	if conf.MQTTBroker != "" {
		mqtt.NewSubscriber(conf.MQTTBroker, conf.MQTTTopic, reportCh)
	}

	reportPeriod := 30
	var err error
	if conf.ReportPeriod != "" {
		reportPeriod, err = strconv.Atoi(conf.ReportPeriod)
		if err != nil {
			fmt.Printf("Error converting %s to integer, period set to default (30)", conf.ReportPeriod)
		}
		fmt.Printf("Report period set to %v minutes\n", conf.ReportPeriod)
	}

	t := transport.NewUDPClient(conf.Master, reportCh, time.Duration(reportPeriod)*time.Minute)
	t.Start()

	pingerProxyPort := 8090
	if conf.PingerProxyPort != "" {
		pingerProxyPort, err = strconv.Atoi(conf.PingerProxyPort)
		if err != nil {
			fmt.Printf("Error converting %s to integer, pinger proxy port set to default (8090)", conf.ReportPeriod)
		}
		fmt.Printf("Pinger proxy port set to: %v\n", conf.PingerProxyPort)
	}
	server := grpc.NewPingerProxy(pingerProxyPort)
	server.Start()

	quitCh := make(chan int)
	<-quitCh
}
