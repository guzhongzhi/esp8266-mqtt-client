#include <Arduino.h>
#include "var.h"
#include "wifi.h"
#include "global.h"
#include "mqtt.h"
#include <ESP8266WiFi.h>

extern String clientId;
extern String userTopic;
extern String MQTTServer;
extern WiFiClient wifiClient;
extern PubSubClient* mqttClient;
int lastMsg = 0;
void setup(){
    clientId = AppId + "-" + String(random(0xffff), HEX);
    Serial.begin(115200);
    Serial.println(clientId);

    pinMode(RelayPin, OUTPUT);
    digitalWrite(RelayPin,LOW);
    RelayStatus = "off";

    if (!autoConfig()){
      smartConfig();
    }
    //irsend.begin();
    userTopic = "/" + AppId + "/user/" +  WiFi.macAddress();
    mqttClient = new PubSubClient(MQTTServer.c_str(),1883,callback,wifiClient);
}

void loop() {

  if (!mqttClient->connected()) {
    mqttReconnect();
  }
  mqttClient->loop();
  long now = millis();
  if ( lastMsg == 0 ||  (now - lastMsg) > 30000) {
    lastMsg = now;
    heartBeat();
  }
}