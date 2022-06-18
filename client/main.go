package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"

	"github.com/jtonynet/autogo/config"
	"github.com/jtonynet/autogo/infrastructure"
)

var addr = flag.String("addr", "localhost:8082", "http service address")
var messageBroker *infrastructure.MessageBroker
var moveTopic string

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		////
		m := string(message)
		fmt.Println("MESSAGE 2: ", m)
		fmt.Println(moveTopic, m)
		messageBroker.Pub(moveTopic, m)
		////

		err = c.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("write:", err)
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	cfg, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	fmt.Println(cfg.MessageBroker.Host)
	messageBroker = infrastructure.NewMessageBroker(cfg.MessageBroker)
	moveTopic = fmt.Sprintf("%s/%s/move", cfg.ProjectName, cfg.RobotName)

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)

	static := http.FileServer(http.Dir("./static"))
	http.Handle("/", static)

	http.ListenAndServe(*addr, nil)

	messageBroker.Sub("autogo/tank-01/sonar", defaultReceiver)
}

func defaultReceiver(client mqtt.Client, msg mqtt.Message) {
	msg.Ack()
	input := "defaultReceiver(\"DEFAULT\" \"" + string(msg.Payload()) + "\")"
	output := "message id:" + strconv.Itoa(int(msg.MessageID())) + " message = " + string(msg.Payload())
	fmt.Println("\n++++++++++++++++")
	fmt.Println(input)
	fmt.Println(output)
	fmt.Println("\n++++++++++++++++")
}

/*
func main2() {
	static := http.FileServer(http.Dir("./static"))
	http.Handle("/", static)

	http.ListenAndServe(":8082", nil)
}
*/
