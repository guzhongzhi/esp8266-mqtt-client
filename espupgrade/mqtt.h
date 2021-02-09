#ifndef MQTT
#define MQTT
void callback(char* topic, byte* payload, unsigned int length);
void mqttReconnect();
void heartBeat();
#endif