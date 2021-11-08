/*
//TODO: Remove. For local test porpouses
package main
*/

package infrastructure

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	config "github.com/jtonynet/autogo/config"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost %v", err)
}

type MQTT struct {
	Client mqtt.Client
}

func NewMQTTClient(cfg config.MQTT) *MQTT {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", cfg.Host, cfg.Port))
	opts.SetClientID("go_mqtt_client")

	//opts.SetUserName("autoGo")
	//opts.SetPassword("******")

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	this := &MQTT{Client: client}
	return this
}

func (this *MQTT) Disconnect(ttl uint) {
	this.Client.Disconnect(ttl)
}

func (this *MQTT) Pub(topic string) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Show msg %d", i)
		token := this.Client.Publish(topic, 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func (this *MQTT) Sub(topic string) {
	token := this.Client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}

/*
//TODO: Remove. For local test porpouses
func main() {
	cfg, err := config.LoadConfig("../.")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	MessageBroker := NewMQTTClient(cfg.MQTT)

	topic := "topic/test"
	MessageBroker.Sub(topic)
	MessageBroker.Pub(topic)

	MessageBroker.Disconnect(250)
}
*/
