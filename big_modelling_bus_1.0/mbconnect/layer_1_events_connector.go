/*
 *
 * Package: mbconnect
 * Layer:   1
 * Module:  events_connector
 *
 * ..... ... .. .
 *
 * Creator: Henderik A. Proper (e.proper@acm.org), TU Wien, Austria
 *
 * Version of: XX.11.2025
 *
 */

package mbconnect

import (
	"github.com/eclipse/paho.mqtt.golang"
	"time"
)

const (
	maxMessageSizeDefault = 10240
)

type (
	tModellingBusEventsConnector struct {
		agentID,
		user,
		port,
		topicRoot,
		broker,
		password string

		messages map[string][]byte

		client mqtt.Client

		reporter *TReporter
	}
)

func (e *tModellingBusEventsConnector) connectionLostHandler(c mqtt.Client, err error) {
	e.reporter.Panic("MQTT connection lost. %s", err)
}

func (e *tModellingBusEventsConnector) connectToMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + e.broker + ":" + e.port)
	opts.SetUsername(e.user)
	opts.SetPassword(e.password)
	opts.SetConnectionLostHandler(e.connectionLostHandler)

	// Apparently not needed:
	//   opts.SetClientID("mqtt-client-" + e.agentID)

	connected := false
	for !connected {
		e.reporter.Progress("Trying to connect to the MQTT broker.")

		e.client = mqtt.NewClient(opts)
		token := e.client.Connect()
		token.Wait()

		err := token.Error()
		if err != nil {
			e.reporter.Error("Error connecting to the MQTT broker. %s", err)

			time.Sleep(5)
		} else {
			connected = true
		}
	}

	e.messages = map[string][]byte{}
	if connected {
		e.reporter.Progress("Connected to the MQTT broker.")

		// Continuously connect all used topics underneath the topic root, and their messages
		// We need this to enable deletion of topics
		mqttTopicPath := e.topicRoot + "/#"
		token := e.client.Subscribe(mqttTopicPath, 1, func(client mqtt.Client, msg mqtt.Message) {
			e.messages[msg.Topic()] = msg.Payload()
		})
		token.Wait()
	}
}

func (e *tModellingBusEventsConnector) listenForEvents(agentID, topicPath string, eventHandler func([]byte)) {
	mqttTopicPath := e.topicRoot + "/" + agentID + "/" + topicPath
	token := e.client.Subscribe(mqttTopicPath, 1, func(client mqtt.Client, msg mqtt.Message) {
		eventHandler(msg.Payload())
	})
	token.Wait()
}

func (e *tModellingBusEventsConnector) messageFromEvent(agentID, topicPath string) []byte {
	mqttTopicPath := e.topicRoot + "/" + agentID + "/" + topicPath
	return e.messages[mqttTopicPath]
}

func (e *tModellingBusEventsConnector) postEvent(topicPath string, message []byte) {
	mqttTopicPath := e.topicRoot + "/" + e.agentID + "/" + topicPath
	token := e.client.Publish(mqttTopicPath, 0, true, string(message))
	token.Wait()
}

func (e *tModellingBusEventsConnector) deleteEvent(topicPath string) {
	e.postEvent(topicPath, []byte{})
}

func createModellingBusEventsConnector(topicBase, agentID string, configData *TConfigData, reporter *TReporter) *tModellingBusEventsConnector {
	e := tModellingBusEventsConnector{}

	e.reporter = reporter

	// Get data from the config file
	e.agentID = agentID
	e.port = configData.GetValue("mqtt", "port").String()
	e.user = configData.GetValue("mqtt", "user").String()
	e.broker = configData.GetValue("mqtt", "broker").String()
	e.password = configData.GetValue("mqtt", "password").String()
	e.topicRoot = configData.GetValue("mqtt", "prefix").String() + "/" + topicBase

	e.connectToMQTT()

	return &e
}
