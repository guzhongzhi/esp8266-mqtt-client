#include <ESP8266WiFi.h>
#include <iostream>
#include <sstream>
#include "ir-send.h"
#include "upgrade.h"
#include "global.h"
#include "ArduinoJson.h"
#include "mqtt.h"
#include <PubSubClient.h>

extern WiFiClient wifiClient;
extern PubSubClient* mqttClient;
extern String MQTTServer;
extern String clientId;
extern String userTopic;
extern String MQTTUser;
extern String MQTTPass;
extern String publicTopic;
extern String RelayStatus;
extern short int RelayPin;
extern short int IRSendPin;
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

  int executedAt = doc["execAt"].as<int>();
  Serial.println("cmd:");
  Serial.println(cmd.c_str());
  
  if ( cmd == "upg" ) {    
    const char* url = doc["data"].as<char*>();
    upgradeUrl = "";
    if(strlen(url) > 0)  {
        upgradeUrl.append(url);
    }
  }
  
  if (cmd == "reset" ) {
    resetFunc();
  }
  /*
   * 
  if(cmd == "ser_send_hex_arr") {
    int len = doc["data"].size();
    delay(500);
    int hexIntData[len];
    for(int i=0;i<len;i++) { 
          const char* v = doc["data"][i].as<char*>();
          int v2 = hex2Int(v);
          hexIntData[i] = v2;
    }
    for (int i=0;i<len;i++){
        Serial.write(hexIntData[i]);
    }
    delay(1200);
    Serial.println("");
  }
  if(cmd == "ser_send_int_arr") {
      int len = doc["data"].size();
      for(int i=0;i<len;i++) {
        double v = data[i];
          Serial.write((int)v);
      }
  }*/
  if(cmd == "srp") {
    uint16_t newRelayPIN = doc["data"].as<uint16_t>();
    if (newRelayPIN != RelayPin) {
      RelayPin = newRelayPIN;
      pinMode(RelayPin, OUTPUT);
      cmd = "off";      
    }
  }
  if( cmd == "irs") {
    const char* data = doc["data"].as<char*>();
    IRSendMessage(IRSendPin,data);
  }
  if( cmd == "pl") {
    int pin = doc["data"].as<int>();
    pinMode(pin, OUTPUT);
    if (pin == RelayPin) {
        RelayStatus = "off";
    }
    digitalWrite(pin,LOW);
  }
  if( cmd == "ph") {
    int pin = doc["data"].as<int>();
    pinMode(pin, OUTPUT);
    digitalWrite(pin,HIGH);
    if (pin == RelayPin) {
        RelayStatus = "on";
    }
  }
  if(cmd == "on") {
    digitalWrite(RelayPin,HIGH);
    RelayStatus = "on";
  }
  if(cmd == "off") {
    digitalWrite(RelayPin,LOW);
    RelayStatus = "off";
  }
  
  mqttClient->publish(heartBeatTopic.c_str(), jsonDeviceInfo(String(cmd.c_str()),executedAt,"feed").c_str());
}

//mqtt 回调
void callback(char* topic, byte* payload, unsigned int length) {
  char data[length + 1];
  for (int i = 0; i < length; i++) {
    data[i] = (char) payload[i];
  }
  data[length] = '\0';
  jsonMessageReceived(data);
}

void mqttReconnect() {
  while (!mqttClient->connected()) {
    Serial.println("mq conn...");
    mqttClient->setBufferSize(2048);
    if (mqttClient->connect(clientId.c_str(),MQTTUser.c_str(),MQTTPass.c_str())) {
      mqttClient->subscribe(userTopic.c_str(),1);
      mqttClient->subscribe(publicTopic.c_str(),1);
      mqttClient->setCallback(callback);
    } else {
      Serial.print("failed, rc=");
      Serial.println(mqttClient->state());
      delay(2000);
    }
  }
}


void heartBeat() {
  if(!mqttClient->connected()) {
    return ;
  }
  mqttClient->publish(heartBeatTopic.c_str(), jsonDeviceInfo("",0,"beat").c_str());
  if(upgradeUrl != "") {
    if(upgrade(upgradeUrl.c_str())) {
      upgradeUrl = "";
    }
  }
}
