package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"hostmonitor/measure"
	"os"
)

type Subscriber struct {
	broker   string
	Client   mqtt.Client
	reportCh chan *measure.Measure
}

func NewSubscriber(broker string, topic string, reportCh chan *measure.Measure) {
	var s Subscriber
	s.broker = broker
	s.reportCh = reportCh
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID("Host-monitor")
	messagePubHandler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Subscriber received a message on topic %s\n", msg.Topic())
		temperature, timestamp := Parse(msg.Payload())
		fmt.Printf("Temperature is %v at %s\n", temperature, timestamp)
		m := &measure.Measure{
			Strings:  make(map[string]string),
			Integers: make(map[string]int64),
			Doubles:  make(map[string]float64),
		}
		m.Integers["MQTT_Temp"] = int64(temperature)
		m.Strings["MQTT_Temp_Time"] = timestamp.String()
		reportCh <- m
	}
	options.SetDefaultPublishHandler(messagePubHandler)
	connectHandler := func(client mqtt.Client) {
		fmt.Println("Subscriber connected")
	}
	options.OnConnect = connectHandler
	connectionLostHandler := func(client mqtt.Client, err error) {
		fmt.Printf("Subscriber connection Lost: %s\n", err.Error())
	}
	options.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
