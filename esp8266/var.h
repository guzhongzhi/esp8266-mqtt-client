#include <PubSubClient.h> //>=2.8
#include <ESP8266WiFi.h>
#include <iostream>
#include <sstream>

PubSubClient *mqttClient;
WiFiClient wifiClient;
bool isNewBoot = true;
String MQTTServer = "118.31.246.195";
String MQTTUser = "mqtt";
String MQTTPass = "mqtt";
String AppId = "camera360";
String clientId = "";
String heartBeatTopic = "/" + AppId + "/heart-beat";
String publicTopic =  "/" + AppId + "/public-topic";
String userTopic =  "";
int RelayPin = 5;
String RelayStatus = "on";
int IRSendPin = 4;
String versionCode = "1.0.0";