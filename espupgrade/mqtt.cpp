#include <ESP8266WiFi.h>
#include <iostream>
#include <sstream>
#include "upgrade.h"
#include "global.h"
#include "ArduinoJson.h"
#include "mqtt.h"
#include <PubSubClient.h> //>=2.8

extern WiFiClient wifiClient;
extern PubSubClient* mqttClient;
extern String MQTTServer;
extern String clientId;
extern String userTopic;
extern String MQTTUser;
extern String MQTTPass;
extern String publicTopic;
extern String RelayStatus;
extern String AppId;
extern bool isNewBoot;
extern String versionCode;
extern String MQTTServer;
extern String heartBeatTopic;
extern string upgradeUrl;

void(* resetFunc) (void) = 0;


void jsonMessageReceived(char* data) {
  StaticJsonDocument<300> doc;
  DeserializationError error = deserializeJson(doc, data);
  Serial.println(error.c_str());
  const char* cmdArray = doc["cmd"].as<char*>();

  string cmd = "";
  cmd.append(cmdArray);

  int executedAt = doc["executedAt"].as<int>();
  Serial.println("cmd:");
  Serial.println(cmd.c_str());
  
  if ( cmd == "upgrade" ) {    
    const char* url = doc["data"].as<char*>();
    upgradeUrl = "";
    if(strlen(url) > 0)  {
        upgradeUrl.append(url);
    }
  }
  
  mqttClient->publish(heartBeatTopic.c_str(), jsonDeviceInfo(String(cmd.c_str()),executedAt,"feedBack").c_str());
}

//mqtt 回调
void callback(char* topic, byte* payload, unsigned int length) {
  char data[length + 1];
  for (int i = 0; i < length; i++) {
    data[i] = (char) payload[i];
  }
  data[length] = '\0';
  Serial.println(data);
  jsonMessageReceived(data);
}

void mqttReconnect() {
  // Loop until we're reconnected
  while (!mqttClient->connected()) {
    Serial.print("mq reconnection...");
    // Attempt to connect
    mqttClient->setBufferSize(2048);
    if (mqttClient->connect(clientId.c_str(),MQTTUser.c_str(),MQTTPass.c_str())) {
      mqttClient->subscribe(userTopic.c_str(),1);
      mqttClient->subscribe(publicTopic.c_str(),1);
      mqttClient->setCallback(callback);
    } else {
      Serial.print("failed, rc=");
      Serial.print(mqttClient->state());
      // Wait 5 seconds before retrying
      delay(2000);
    }
  }
}


void heartBeat() {
  if(!mqttClient->connected()) {
    return ;
  }
  mqttClient->publish(heartBeatTopic.c_str(), jsonDeviceInfo("",0,"heartBeat").c_str());
  if(upgradeUrl != "") {
    if(upgrade(upgradeUrl.c_str())) {
      upgradeUrl = "";
    }
  }
}
