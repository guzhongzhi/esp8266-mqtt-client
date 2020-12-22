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
//红外发射头文件
#include <assert.h>
#include <IRrecv.h>
#include <IRremoteESP8266.h>
#include <IRac.h>
#include <IRtext.h>
#include <IRutils.h>
#include "ArduinoJson.h"

StaticJsonDocument<300> doc;

using namespace std;

WiFiClient espClient;

const bool JSONEnabled = true; //是否使用JSON通信

//继电器及状态LED
const uint16_t statePIN = 2;  //ESP8266 GPIO pin to use. Recommended: 2 . 开机状态
const uint16_t relayPIN = 0; //ESP8266 GPIO pin to use. Recommended: 0  继电器
String relayPINState = "off";

//红外发射
const uint16_t kIrLed = 4; // ESP8266 GPIO pin to use. Recommended: 4 (D2). 红外
IRsend irsend(kIrLed);     // Set the GPIO to be used to sending the message.

//MQTT
String APP_ID = "guz";
String clientId = "";
unsigned long lastMsg = 0;
//String MQTT_SERVER = "118.31.246.195";
String MQTT_SERVER = "192.168.18.159";

//红外接收
int isIrEnabled = 1; //是否启用红外输入
const uint16_t kRecvPin = 2;
const uint16_t kCaptureBufferSize = 1024;
#if DECODE_AC
const uint8_t kTimeout = 50;
#else   // DECODE_AC
const uint8_t kTimeout = 15;
#endif  // DECODE_AC
const uint16_t kMinUnknownSize = 12;
const uint8_t kTolerancePercentage = kTolerance;  // kTolerance is normally 25%
#define LEGACY_TIMING_INFO false
IRrecv irrecv(kRecvPin, kCaptureBufferSize, kTimeout, true);
decode_results results;  // Somewhere to store the results
    
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
    Serial.println( WiFi.gatewayIP());
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
    Serial.println( WiFi.gatewayIP());
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
    Serial.println( WiFi.gatewayIP());
}

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived [");
  Serial.print(topic);
  Serial.print(",length:");
  Serial.print(length);
  Serial.print("] ");
  
  char data[length + 1];
  for (int i = 0; i < length; i++) {
    data[i] = (char) payload[i];
  }
  data[length] = '\0';
  Serial.println("");
  Serial.print(data);

  if(JSONEnabled) {
    jsonMessageReceived(data);
    return;
  }

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

void jsonMessageReceived(char* data) {
  StaticJsonDocument<300> doc;
  DeserializationError error = deserializeJson(doc, data);
  Serial.println("JSONDecode error");
  Serial.println(error.c_str());
  const char* cmd = doc["cmd"].as<char*>();
  Serial.println("");
  Serial.print("cmd:");
  Serial.print(cmd);
  Serial.println("");
  if(cmd == "irs" || cmd == "irSend") {
    const char* data = doc["cmd"].as<char*>();
    sendCode(data,"");
  }
  if(cmd == "setPinLow") {
    int pin = doc["pin"].as<int>();
    digitalWrite(pin,LOW);
  }
  if(cmd == "setPinHigh") {
    int pin = doc["pin"].as<int>();
    digitalWrite(pin,HIGH);
  }
}

PubSubClient client(MQTT_SERVER.c_str(),1883,callback,espClient);

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT "+MQTT_SERVER+" connection...");
    // Attempt to connect
    client.setBufferSize(2048);
    String publicTopic =  "/" + APP_ID + "/public-topic";
    if (client.connect(clientId.c_str(),"mqtt","mqtt")) {
      Serial.println("connected");
      String ss =  ("/" + APP_ID + "/user/" +  WiFi.macAddress());
      client.subscribe(ss.c_str(),1);
      client.subscribe(publicTopic.c_str(),1);
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
    relayPINState = "on";
    Serial.println("replay high");
    digitalWrite(relayPIN,HIGH);
    heartBeat();
}

void setLow() {
  Serial.println("replay low");
  relayPINState = "off";
  //digitalWrite(relayPIN,LOW);
  analogWrite(relayPIN,0);
  heartBeat();
}

void heartBeat() {
  if(!client.connected()) {
    return ;
  }
  delay(500);
  String heartBeatTopic = "/" + APP_ID + "/heart-beat";
  Serial.println(heartBeatTopic);
  
  if(JSONEnabled) {
    client.publish(heartBeatTopic.c_str(), jsonDeviceInfo("").c_str());
  } else {
    Serial.println(deviceInfo().c_str());
    client.publish(heartBeatTopic.c_str(), deviceInfo().c_str());
  }
}

void irReceived(String data) {
  if(!isIrEnabled) {
      return;
  }
  String commonInfo = "";
    if(JSONEnabled) {
        commonInfo = jsonDeviceInfo(data);
    } else {
      commonInfo = deviceInfo();
      if(data.length() > 0) {
        commonInfo += ("&data=" + data);
      }  
    }
    client.publish(("/" + APP_ID + "/ir-received").c_str(), commonInfo.c_str());
}

String jsonDeviceInfo(String data) {
   doc["mac"] = WiFi.macAddress();
   doc["ip"]   = WiFi.localIP().toString();
   doc["jsonEnabled"] = JSONEnabled;
   doc["wifi"] = WiFi.SSID();
   doc["clientId"] = clientId;
   doc["gw"] = WiFi.gatewayIP().toString();
   doc["relay"] = relayPINState.c_str();
   doc["relayPin"] = relayPIN;
   doc["statePin"] = statePIN;
   doc["irPin"] = kIrLed;
   doc["appName"] = APP_ID;
   doc["data"] = data;
   String output = "";
   serializeJson( doc,  output);
   Serial.println(output);
   return output;
}

String deviceInfo() {
  String s = "mac=";
  s.concat(WiFi.macAddress());
  s.concat("&ip=");
  s.concat(WiFi.localIP().toString());
  s.concat("&wifi=");
  s.concat(WiFi.SSID().c_str());
  s.concat("&clientId=");
  s.concat(clientId);
  s.concat("&gw=");
  s.concat(WiFi.gatewayIP().toString());
  s.concat("&relay=");
  s.concat(relayPINState.c_str());
  return s;
}

void setup(void)
{
  clientId = APP_ID + "-" + String(random(0xffff), HEX);
  Serial.begin(115200);
  Serial.println("");
  pinMode(relayPIN, OUTPUT);
  pinMode(statePIN, OUTPUT);
  setLow();
  int DEBUG = 0;
  if(DEBUG == 1) {
    debugWIFI();
  } else if (!autoConfig()){
      Serial.println( "Start AP mode" );
      smartConfig();
  }
  delay(2000);
  if(isIrEnabled == 1 ) {
    irsend.begin();
    setupIR();    
  }
}

void loop(void)
{
  if (!client.connected()) {
    reconnect();
  }
  client.loop();
  long now = millis();
  if ( lastMsg == 0 ||  (now - lastMsg) > 10000) {
    lastMsg = now;
    heartBeat();
  }
  if(isIrEnabled == 1) {
    checkIrInput();
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





//红外接收

void setupIR() {
  assert(irutils::lowLevelSanityCheck() == 0);
  Serial.printf("\n" D_STR_IRRECVDUMP_STARTUP "\n", kRecvPin);
#if DECODE_HASH
  irrecv.setUnknownThreshold(kMinUnknownSize);
#endif  // DECODE_HASH
  irrecv.setTolerance(kTolerancePercentage);  // Override the default tolerance.
  irrecv.enableIRIn();  // Start the receiver
}

void checkIrInput() {
  // Check if the IR code has been received.
  if (irrecv.decode(&results)) {
    // Display a crude timestamp.
    uint32_t now = millis();
    Serial.printf(D_STR_TIMESTAMP " : %06u.%03u\n", now / 1000, now % 1000);
    // Check if we got an IR message that was to big for our capture buffer.
    if (results.overflow)
      Serial.printf(D_WARN_BUFFERFULL "\n", kCaptureBufferSize);
    // Display the library version the message was captured with.
    Serial.println(D_STR_LIBRARY "   : v" _IRREMOTEESP8266_VERSION_ "\n");
    // Display the tolerance percentage if it has been change from the default.
    if (kTolerancePercentage != kTolerance)
      Serial.printf(D_STR_TOLERANCE " : %d%%\n", kTolerancePercentage);
    // Display the basic output of what we found.
    Serial.print(resultToHumanReadableBasic(&results));
    // Display any extra A/C info if we have it.
    String description = IRAcUtils::resultAcToString(&results);
    if (description.length()) Serial.println(D_STR_MESGDESC ": " + description);
    yield();  // Feed the WDT as the text output can take a while to print.
#if LEGACY_TIMING_INFO
    // Output legacy RAW timing info of the result.
    Serial.println(resultToTimingInfo(&results));
    yield();  // Feed the WDT (again)
#endif  // LEGACY_TIMING_INFO
    // Output the results as source code
    String b = formatIRData(resultToSourceCode(&results));
    Serial.println("=======================");
    Serial.println(b);
    Serial.println("=======================");
    irReceived(b);
    Serial.println();    // Blank line between entries
    yield();             // Feed the WDT (again)
  }
}

String formatIRData(String m) {
    String n = "";
    int isStarted = 0;
    for(int i=0;i<m.length();i++) {
       if(m[i] == ' ')  {
          continue;
       }
       if(m[i] == '\n' || m[i]=='\r') {
           break;
       }
       if(m[i] == '{') {
        isStarted = 1;
        continue;
       }
       if(m[i]=='}') {
        break;
       }
       if(isStarted) {
         n += String(m[i]);
       }
    }   
    Serial.println(n);
    return n;
}
