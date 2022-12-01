package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"hostmonitor/grpc"
	"hostmonitor/measure"
	"hostmonitor/probers"
	"hostmonitor/sensors"
	"hostmonitor/transport"
	"io/ioutil"
	"log"
	"net"
	"os"
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
}

func loadConfiguration(path string) *Configuration {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Could not open %s error: %s\n", path, err)
		conf := &Configuration{true, true, true, true, true, true, "127.0.0.1:8758", "127.0.0.1", "30", "8090"}
		log.Printf("Host Monitor will use default configuration: %v\n", conf)
		return conf
	}
	if yfile == nil {
		panic("There was no error but YFile was null")
	}
	conf := Configuration{Master: "127.0.0.1:8758"}
	err2 := yaml.Unmarshal(yfile, &conf)
	if err2 != nil {
		log.Printf("Configuration file could not be parsed, error: %s\n", err2)
		panic(err2)
	}

	log.Printf("Found configuration: %v\n", conf)
	return &conf
}

// Get IPv4 from Network Interfaces
func localAddress() string {
	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error retrieving hostname: %v\n", err.Error())
	}
	addrs, err2 := net.LookupIP(host)
	if err2 != nil {
		log.Fatalf("Error retrieving local addresses: %v\n", err2.Error())
	}
	ip := ""
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip = ipv4.String()
		}
	}
	return ip
}

func main() {
	version := "11-27-2022"
	fmt.Println("Running software version ", version)
	file, errl := os.OpenFile("./log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if errl != nil {
		log.Fatalf("Unable to create log file\n%s", errl)
	}
	log.SetOutput(file)
	fmt.Println("Program output will be logged in ./log")

	conf := loadConfiguration("hostmonitor.yaml")
	reportCh := make(chan *measure.Measure)

	var boardAddress string
	if conf.BoardIP != "" {
		boardAddress = conf.BoardIP
	} else {
		outboundIP := localAddress()
		s := strings.Split(outboundIP, ".")
		s[len(s)-1] = "1"
		boardAddress = strings.Join(s[:], ".")
		log.Println("The board address was not specified so it was automatically detected as: " + boardAddress)
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

	reportPeriod := 30
	var err error
	if conf.ReportPeriod != "" {
		reportPeriod, err = strconv.Atoi(conf.ReportPeriod)
		if err != nil {
			log.Printf("Error converting %s to integer, period set to default (30)", conf.ReportPeriod)
		}
		log.Printf("Report period set to %v minutes\n", conf.ReportPeriod)
	}

	t := transport.NewUDPClient(conf.Master, reportCh, time.Duration(reportPeriod)*time.Minute)
	t.Start()

	pingerProxyPort := 8090
	if conf.PingerProxyPort != "" {
		pingerProxyPort, err = strconv.Atoi(conf.PingerProxyPort)
		if err != nil {
			log.Printf("Error converting %s to integer, pinger proxy port set to default (8090)", conf.ReportPeriod)
		}
		log.Printf("Pinger proxy port set to: %v\n", conf.PingerProxyPort)
	}
	server := grpc.NewPingerProxy(pingerProxyPort)
	server.Start()

	// Flush the log every week
	go func() {
		logFlush := time.NewTicker(time.Hour * 24 * 7)
		for {
			select {
			case <-logFlush.C:
				if err != nil {
					log.Fatalf("Unable to clear log file\n%s", err)
				}
			}
		}
	}()

	quitCh := make(chan int)
	<-quitCh
}
