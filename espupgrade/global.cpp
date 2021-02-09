#include <iostream>
#include <sstream>
#include <ESP8266WiFi.h>
#include "ArduinoJson.h"
#include "global.h"
using namespace std;

extern String clientId;
extern String AppId;
extern bool isNewBoot;
extern String versionCode;

String jsonDeviceInfo(String data, int executedAt, String cmd) {
   StaticJsonDocument<400> doc;
   doc["mac"] = WiFi.macAddress();
   doc["ip"]   = WiFi.localIP().toString();
   doc["wifi"] = WiFi.SSID();
   doc["clientId"] = clientId;
   doc["gw"] = WiFi.gatewayIP().toString();
   doc["relay"] = 0;
   doc["relayPin"] = 5;
   doc["statePin"] = -1;
   doc["irPin"] = 4;
   doc["appName"] = AppId;
   doc["data"] = data;
   doc["cmd"] = cmd;
   doc["isNewBoot"] = isNewBoot;
   doc["executedAt"] = executedAt;
   doc["version"] = versionCode;
   String output = "";
   serializeJson( doc,  output);
   Serial.println(output);
   isNewBoot = false;
   return output;
}
