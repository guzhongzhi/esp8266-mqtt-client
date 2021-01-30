#include <iostream>
#include <sstream>
#include <ESP8266WiFi.h>
#include "ArduinoJson.h"
#include "global.h"
using namespace std;

string replaceCommaToSpace(string s) {
  int n = s.length();
  for (int i = 0; i < n; ++i){
    if (s[i] == ','){
      s[i] = ' ';
    }
  }
  return s;
}


//hex string convert to int
int hex2Int(string v)  {
  int temp;
  std::stringstream ss2;
  ss2 << std::hex <<v;
  ss2 >> temp;
  return temp;
}

extern String clientId;
extern String RelayStatus;
extern short int RelayPin;
extern short int IRSendPin;
extern String AppId;
extern bool isNewBoot;
extern String versionCode;

String jsonDeviceInfo(String data, int executedAt, String cmd) {
   StaticJsonDocument<400> doc;
   doc["mac"] = WiFi.macAddress();
   doc["ip"]   = WiFi.localIP().toString();
   doc["wifi"] = WiFi.SSID();
   doc["cid"] = clientId;
   doc["gw"] = WiFi.gatewayIP().toString();
   doc["relay"] = RelayStatus.c_str();
   doc["rPin"] = RelayPin;
   doc["sPin"] = -1;
   doc["irPin"] = IRSendPin;
   doc["app"] = AppId;
   doc["data"] = data;
   doc["cmd"] = cmd;
   doc["newBoot"] = isNewBoot;
   doc["execAt"] = executedAt;
   doc["ver"] = versionCode;
   String output = "";
   serializeJson( doc,  output);
   Serial.println(output);
   isNewBoot = false;
   return output;
}
