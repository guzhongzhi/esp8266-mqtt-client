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
#define MQTT_MAX_PACKET_SIZE 2048
#ifndef UNIT_TEST
#include <Arduino.h>
#endif
#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <PubSubClient.h> //>=2.8
#include<string>
using namespace std;
const char *ssid = "10012503";
const char *password = "gd10012503";

WiFiClient espClient;
unsigned long lastMsg = 0;
#define MSG_BUFFER_SIZE  (5000)
char msg[MSG_BUFFER_SIZE];
int value = 0;

bool autoConfig()
{
    int a = 0;
    WiFi.begin();
        while ( WiFi.status() != WL_CONNECTED )
        {
            Serial.println( "AutoConfig Success" );
            Serial.printf( "SSID:%s\r\n", WiFi.SSID().c_str() );
            Serial.printf( "PSW:%s\r\n", WiFi.psk().c_str() );
            WiFi.printDiag( Serial );
            delay( 1000 );
            a++;
            if ( a == 10 )
            {
                a = 0;
                return(false);
                break;
            }
        }
        if ( false )
        {
            Serial.println( "" );
            Serial.println( "wifi line faild !" );
        }else  {
            Serial.println( "" );
            Serial.println( "WiFi connected" );
            Serial.println( "IP address: " );
            Serial.println( WiFi.localIP() );
            return(true);
        }
}
void smartConfig()
{
    WiFi.mode( WIFI_STA );
    Serial.println( "\r\nWait for Smartconfig" );
    WiFi.beginSmartConfig();
    while ( 1 )
    {
        Serial.print( "Wait soft line..\r\n" );
        if ( WiFi.smartConfigDone() )
        {
            Serial.println( "SmartConfig Success" );
            Serial.printf( "SSID:%s\r\n", WiFi.SSID().c_str() );
            Serial.printf( "PSW:%s\r\n", WiFi.psk().c_str() );
            WiFi.setAutoConnect( true ); /* 设置自动连接 */
            break;
        }
        delay( 1000 );
    }
    Serial.println( "" );
    Serial.println( "WiFi connected" );
    Serial.println( "IP address: " );
    Serial.println( WiFi.localIP() );
}

void debugWIFI() {
  WiFi.begin(ssid, password);
    Serial.println("");
    // Wait for connection
    while (WiFi.status() != WL_CONNECTED)
    {
      delay(500);
      Serial.print(".");
    }
    Serial.println("");
    Serial.print("Connected to ");
    Serial.println(ssid);
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());
}


void callback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived [");
  Serial.println(topic);
  Serial.println(length);
  Serial.print("] ");
  for (int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println("");
  Serial.println("=============================");
}

PubSubClient client("s1.gulusoft.com",1883,callback,espClient);

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Create a random client ID
    String clientId = "ESP8266Client-";
    clientId += String(random(0xffff), HEX);
    // Attempt to connect
    client.setBufferSize(2048);
    if (client.connect(clientId.c_str(),"admin","admin","test",1,false,"")) {
      Serial.println("connected");
      client.subscribe("test",1);
      client.setCallback(callback);
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}

void setup(void)
{
  Serial.begin(115200);
  Serial.println("");

  int DEBUG = 0;
  if(DEBUG == 1) {
    debugWIFI();
  } else {
    if ( !autoConfig() )
    {
        Serial.println( "Start AP mode" );
        smartConfig();
    }
  } 
}

void loop(void)
{
  if (!client.connected()) {
    reconnect();
  }
  client.loop();
}
