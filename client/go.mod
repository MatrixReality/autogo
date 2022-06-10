module client

go 1.16

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/gorilla/websocket v1.5.0
	github.com/jtonynet/autogo/config v0.0.0
	github.com/jtonynet/autogo/infrastructure v0.0.0
)

replace github.com/jtonynet/autogo/config => ../config

replace github.com/jtonynet/autogo/infrastructure => ../infrastructure
