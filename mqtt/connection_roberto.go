package mqtt

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"log"
	"net/url"
	"os"
)

type MQTT5 struct {
	serverURL string
	ctx       context.Context
	cancel    context.CancelFunc
	cm        *autopaho.ConnectionManager
}

func (mqtt *MQTT5) Connect() {
	parsedURL, e := url.Parse(mqtt.serverURL)
	if e != nil {
		log.Fatal("MQTT URL parse failed: ", e)
		return
	}

	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{parsedURL},
		KeepAlive:         30,
		ConnectRetryDelay: 10000,
		OnConnectionUp:    func(*autopaho.ConnectionManager, *paho.Connack) { fmt.Println("mqtt connection up") },
		OnConnectError:    func(err error) { fmt.Println("error whilst attempting connection: ", err) },
		Debug:             paho.NOOPLogger{},
		ClientConfig: paho.ClientConfig{
			ClientID:      "qttBroker",
			OnClientError: func(err error) { fmt.Println("server requested disconnect: ", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					fmt.Println("server requested disconnect: ", d.Properties.ReasonString)
				} else {
					fmt.Println("server requested disconnect; reason code: ", d.ReasonCode)
				}
			},
		},
	}
	mqtt.ctx, mqtt.cancel = context.WithCancel(context.Background())
	var err error
	mqtt.cm, err = autopaho.NewConnection(mqtt.ctx, cliCfg)
	if err != nil {
		fmt.Println("Connection failed ", err)
		os.Exit(-1)
	}
}
