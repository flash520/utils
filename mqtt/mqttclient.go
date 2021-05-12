/**
 * @Author: koulei
 * @Description: TODO
 * @File:  namingClient
 * @Version: 1.0.0
 * @Date: 2021/5/12 00:48
 */

package mqtt

import (
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"sync"
	"time"
)

var mqttClient mqtt.Client

type MqttClient struct {
}

// CreateMqttClient 初始化mqtt客户端, addr: 服务端地址和端口, clientID,username 授权的客户端ID和用户名
func CreateMqttClient(addr, clientID, username string) *MqttClient {
	var once sync.Once
	once.Do(func() {
		//mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
		opts := mqtt.NewClientOptions().AddBroker("tcp://" + addr).
			SetClientID(clientID).SetUsername(username)

		opts.SetKeepAlive(600 * time.Second)
		//opts.SetDefaultPublishHandler(f)
		opts.SetPingTimeout(1 * time.Second)
		opts.SetAutoReconnect(true)
		opts.SetConnectRetry(true)

		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		mqttClient.Disconnect(250)
	})

	return &MqttClient{}
}

// GetConn 获取一个mqtt连接
func (e *MqttClient) GetConn() mqtt.Client {
	return mqttClient
}
