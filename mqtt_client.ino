/*
 * Common mistakes & tips:
 *   * Don't just connect the IR LED directly to the pin, it won't
 *     have enough current to drive the IR LED effectively.
 *   * Make sure you have the IR LED polarity correct.
 *     See: https://learn.sparkfun.com/tutorials/polarity/diode-and-led-polarity
 *   * Typical digital camera/phones can be used to see if the IR LED is flashed.
 *     Replace the IR LED with a normal LED if you don't have a digital camera
 *     when debugging.
 *   * Avoid using the following pins unless you really know what you are doing:
 *     * Pin 0/D3: Can interfere with the boot/program mode & support circuits.
 *     * Pin 1/TX/TXD0: Any serial transmissions from the ESP8266 will interfere.
 *     * Pin 3/RX/RXD0: Any serial transmissions to the ESP8266 will interfere.
 *   * ESP-01 modules are tricky. We suggest you use a module with more GPIOs
 *     for your first time. e.g. ESP-12 etc.
 *     https://github.com/abhrodeep/Arduino_projs
 */

#ifndef UNIT_TEST
#include <Arduino.h>
#endif 
#define MQTT_MAX_PACKET_SIZE 2048
#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <PubSubClient.h> //>=2.8
#include <IRremoteESP8266.h>
#include <IRsend.h>
#include <vector>
#include <iostream>
#include <sstream> 
#include<string>
#include "ESP8266HTTPClient.h"

using namespace std;

WiFiClient espClient;
unsigned long lastMsg = 0;
#define MSG_BUFFER_SIZE  (5000)
char msg[MSG_BUFFER_SIZE];
int value = 0;

const uint16_t statePIN = 0;  //ESP8266 GPIO pin to use. Recommended: 0 (D3). 开机状态
const uint16_t relayPIN = 5; //ESP8266 GPIO pin to use. Recommended: 5 (D1). 继电器
const uint16_t kIrLed = 4; // ESP8266 GPIO pin to use. Recommended: 4 (D2). 红外
IRsend irsend(kIrLed);     // Set the GPIO to be used to sending the message.

bool autoConfig()
{
    int tried = 0;
    Serial.println( "" );
    Serial.print( "Start to connect WIFI." );
    WiFi.begin();
    while ( WiFi.status() != WL_CONNECTED )
    {
        Serial.print( "." );
        tried++;
        delay( 1000 );
        if(tried >= 20) {
          Serial.println( "" );
          return false;
        }
    }
    Serial.println( "" );
    Serial.println( "WiFi connected" );
    Serial.println( "IP address: " );
    Serial.println( WiFi.localIP() );
    return(true);
}
void smartConfig()
{
    Serial.println( "\r\nWait for Smartconfig" );
    WiFi.mode( WIFI_STA );
    WiFi.beginSmartConfig();
    Serial.print( "Wait soft line.." );
    while ( WiFi.status() != WL_CONNECTED )
    {
        Serial.print(".");
        if (! WiFi.smartConfigDone() )
        {
            delay( 1000 );
            continue;
        }
        Serial.println( "SmartConfig Success" );
        Serial.printf( "SSID:%s\r\n", WiFi.SSID().c_str() );
        Serial.printf( "PSW:%s\r\n", WiFi.psk().c_str() );
        WiFi.setAutoConnect( true ); /* 设置自动连接 */
        break;        
    }
    Serial.println( "" );
    Serial.println( "WiFi connected" );
    Serial.println( "IP address: " );
    Serial.println( WiFi.localIP() );
}

void debugWIFI() {
    WiFi.begin("10012503", "gd10012503");
    Serial.println("");
    // Wait for connection
    while (WiFi.status() != WL_CONNECTED)
    {
      delay(500);
      Serial.print(".");
    }
    Serial.println("");
    Serial.print("Connected to ");
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());
}

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived [");
  Serial.print(topic);
  Serial.print(length);
  Serial.print("] ");
  for (int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println("");

  string cmd = "";
  string message = "";
  int isMessage = 0;
  for(int i =0;i<length;i++) {
    char c = (char)payload[i];
    if(c == '|') {
       isMessage = true;
       continue;
    } else if(isMessage == 0){
      cmd += c;
    } else if(isMessage > 0) {
      message += (char)payload[i];
    }
  }
  if(cmd == "irs") {
    sendCode(message,"");
  } else if(cmd == "upp" || cmd == "on") {
    setHigh();
  } else if(cmd == "low" || cmd == "off") {
    setLow();
  }
  
  Serial.println("");
  Serial.println("=============================");
}

PubSubClient client("192.168.18.60",1883,callback,espClient);
//PubSubClient client("s1.gulusoft.com",1883,callback,espClient);

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Create a random client ID
    String clientId = "addf59cad3fb9-";
    clientId += String(random(0xffff), HEX);
    // Attempt to connect
    client.setBufferSize(2048);
    char * topic = "addf59cad3fb9-topic";
    char * globalTopic = "addf59cad3fb9-global";
    if (client.connect(clientId.c_str(),"admin","admin",topic,1,false,"")) {
      Serial.println("connected");
      client.subscribe(topic,1);
      client.subscribe(globalTopic,1);
      client.setCallback(callback);
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 2 seconds");
      // Wait 5 seconds before retrying
      delay(2000);
    }
  }
}

void setHigh() {
    Serial.println("replay high");
    digitalWrite(relayPIN,HIGH);
}

void setLow() {
  Serial.println("replay low");
  digitalWrite(relayPIN,LOW);
}


void sendHttpOut(String data) {
    HTTPClient http;

    String s1 = "http://esp8266.gulusoft.com/index.php?mac=";
    s1.concat(WiFi.macAddress());
    s1.concat("&ip=");
    s1.concat(WiFi.localIP().toString());
    s1.concat("&wifi=");
    s1.concat(WiFi.SSID().c_str());
    s1.concat("&data=");
    s1.concat(data);
    Serial.println(s1);
    http.begin(s1); 
    http.addHeader("Content-Type", "application/json"); 
    http.GET();
    http.end();
}

void setup(void)
{
  Serial.begin(115200);
  Serial.println("");
  pinMode(relayPIN, OUTPUT);
  pinMode(statePIN, OUTPUT);
  digitalWrite(statePIN,HIGH);
  digitalWrite(relayPIN,HIGH);
  int DEBUG = 0;
  if(DEBUG == 1) {
    debugWIFI();
  } else if (!autoConfig()){
      Serial.println( "Start AP mode" );
      smartConfig();
  }
  delay(2000);
  sendHttpOut("init");
  irsend.begin();
}

void loop(void)
{
  if (!client.connected()) {
    reconnect();
  }
  client.loop();
  long now = millis();
  if (now - lastMsg > 2000) {
    lastMsg = now;
    client.publish("addf59cad3fb9-hart-beat", WiFi.macAddress().c_str());
  }
}

string replaceCommaToSpace(string s) {
  int n = s.length();
  for (int i = 0; i < n; ++i){
    if (s[i] == ','){
      s[i] = ' ';
    }
  }
  return s;
}

//
void sendCode(string message, string type) {
  message = replaceCommaToSpace(message);
  istringstream is(message);
  vector<uint16_t> v;
  uint16_t i;
  while(is>>i)
  {
      v.push_back(i);
  }
  uint16_t rawData[v.size()];
  for(int i=0;i<v.size();i++) {
    rawData[i] = v[i];
  }
  
  Serial.println("start to send IR");
  irsend.sendRaw(rawData, v.size(), 38);
  Serial.println("end to send IR");
}
