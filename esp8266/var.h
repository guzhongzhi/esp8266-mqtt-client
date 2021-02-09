#include <PubSubClient.h>
#include <ESP8266WiFi.h>
#include <iostream>
#include <sstream>
#include <IRrecv.h>

using namespace std;

PubSubClient *mqttClient;
WiFiClient wifiClient;
bool isNewBoot = true;
String MQTTServer = "118.31.246.195";
String MQTTUser = "mqtt";
String MQTTPass = "mqtt";
String AppId = "guz";
String heartBeatTopic = "/" + AppId + "/heart-beat";
String publicTopic =  "/" + AppId + "/public-topic";
String userTopic =  "";
short int RelayPin = 4;
String RelayStatus = "off";
short int IRSendPin = 14;
short int IRReceivePin = 5;
String versionCode = "1.2.4";
bool isInUpgrading = false;
string upgradeUrl = "";
decode_results results;
