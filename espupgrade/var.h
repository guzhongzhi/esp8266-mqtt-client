#include <PubSubClient.h> //>=2.8
#include <ESP8266WiFi.h>
#include <iostream>
#include <sstream>

using namespace std;

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
String versionCode = "1.2.0_upgrader";
bool isInUpgrading = false;
string upgradeUrl = "";
