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
        WiFi.setAutoConnect( true ); 
        break;
    }
}

bool autoConfig()
{
    Serial.println( "connect wifi." );
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
    return(WiFi.status() == WL_CONNECTED);
}
