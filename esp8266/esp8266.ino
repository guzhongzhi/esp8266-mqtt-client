#include <Arduino.h>
#include "var.h"
#include "wifi.h"
#include "global.h"
#include "mqtt.h"
#include <ESP8266WiFi.h>

//红外发射头/接收文件
#include <IRac.h>

extern String userTopic;
extern String MQTTServer;
extern String AppId;
extern WiFiClient wifiClient;
extern PubSubClient* mqttClient;
extern short int IRReceivePin;
int lastMsg = 0;
extern bool isInUpgrading;
extern bool isNewBoot;
extern decode_results results;

#define LEGACY_TIMING_INFO false
const uint8_t kTolerancePercentage = kTolerance;  // kTolerance is normally 25%
IRrecv irrecv(IRReceivePin, 1024, 15, true);


void setupIRReceive() {
  irrecv.setTolerance(15);  // Override the default tolerance.
  irrecv.enableIRIn();      // Start the receiver
}
void setup(){
    Serial.begin(115200);
    delay(1000);
    pinMode(RelayPin, OUTPUT);
    digitalWrite(RelayPin,LOW);
    if (!autoConfig()){
      smartConfig();
    }
    userTopic = "/" + AppId + "/user/" +  WiFi.macAddress();
    mqttClient = new PubSubClient(MQTTServer.c_str(),1883,callback,wifiClient);
    setupIRReceive();
}


bool checkIrInput() {
  if (irrecv.decode(&results)) {
    String a = resultToSourceCode(&results);
    String b = jsonDeviceInfo(formatIRData2(a),0,"irr");
    mqttClient->publish(heartBeatTopic.c_str(), b.c_str());
    delay(300);
    return true;
  }
  return false;
}

String formatIRData2(String m) {
    String n = "{";
    int isStarted = 0;
    for(int i=0;i<m.length();i++) {
       if(m[i] == '{' ) {
         isStarted = 1;
         continue;
       }
       if(isStarted != 1 || m[i] == ' ') {
        continue;
       }
       if(m[i] == '\r' || m[i] == '\n') {
         break;
       }
        n += String(m[i]);      
    }
    return n;
}

void loop() {
  if (!mqttClient->connected()) {
    mqttReconnect();
  }
  mqttClient->loop();
  long now = millis();
  if(checkIrInput()) {
    lastMsg = now;
    return;
  }
  if ( lastMsg == 0 ||  (now - lastMsg) > 10000) {
    lastMsg = now;
    heartBeat();
  }
}
