#include <Arduino.h>
#include <iostream>
#include <sstream>

using namespace std;

#ifndef GLOBAL
#define GLOBAL
String jsonDeviceInfo(String data, int executedAt, String cmd);
string replaceCommaToSpace(string);
int hex2Int(string v);
#endif