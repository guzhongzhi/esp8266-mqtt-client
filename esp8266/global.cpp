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

extern String RelayStatus;
extern short int RelayPin;
extern short int IRSendPin;
extern String AppId;
extern bool isNewBoot;
extern String versionCode;
extern short int IRReceivePin;

String jsonDeviceInfo(String data, int executedAt, String cmd) {
   StaticJsonDocument<600> doc;
   doc["m"] = WiFi.macAddress();
   doc["i"]   = WiFi.localIP().toString();
   doc["w"] = WiFi.SSID();
   //doc["g"] = WiFi.gatewayIP().toString();
   doc["r"] = RelayStatus.c_str();
   doc["rp"] = RelayPin;
   doc["sp"] = -1;
   doc["irsp"] = IRSendPin;
   doc["irrp"] = IRReceivePin;
   doc["a"] = AppId;
   doc["d"] = data;
   doc["c"] = cmd;
   doc["b"] = isNewBoot;
   doc["e"] = executedAt;
   doc["v"] = versionCode;
   String output = "";
   serializeJson( doc,  output);
   isNewBoot = false;
   return output;
}
