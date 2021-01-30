#include <ESP8266WiFi.h>
#include "wifi.h"

void smartConfig()
{
    WiFi.mode( WIFI_STA );
    WiFi.beginSmartConfig();
    while ( WiFi.status() != WL_CONNECTED )
    {
        if (! WiFi.smartConfigDone() )
        {
            delay( 1000 );
            continue;
        }
        Serial.printf( "SSID:%s\r\n", WiFi.SSID().c_str() );
        Serial.printf( "PSW:%s\r\n", WiFi.psk().c_str() );
        WiFi.setAutoConnect( true ); 
        break;
    }
    Serial.println( WiFi.localIP() );
    Serial.println( WiFi.gatewayIP());
}

bool autoConfig()
{
    Serial.println( "Start to connect WIFI." );
    WiFi.begin();
    short int maxNum = 10;
    while ( maxNum  > 0)
    {
        maxNum--;
        Serial.print( "." );
        delay( 1000 );
        if(WiFi.status() != WL_CONNECTED) {
          continue;
        }
    }
    bool i = WiFi.status() == WL_CONNECTED;
    if(i) {
      Serial.println( WiFi.localIP() );
      Serial.println( WiFi.gatewayIP());      
    }
    return(i);
}
