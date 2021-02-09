#include <ESP8266WiFi.h>
#include "wifi.h"

//智能配网
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
}


//自动连接
bool autoConfig()
{
    Serial.println( "Connect WIFI." );
    WiFi.begin();
    while ( WiFi.status() != WL_CONNECTED )
    {
        Serial.print( "." );
        delay( 1000 );
    }
    Serial.println( WiFi.localIP() );
    return(true);
}
