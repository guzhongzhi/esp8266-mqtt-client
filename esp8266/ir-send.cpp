#include <IRsend.h>
#include "ir-send.h"
#include <iostream>
#include <sstream>
#include <vector>
#include "global.h"

using namespace std;

void IRSendMessage(short int pin,string message) {
  IRsend irsend(pin); 
  irsend.begin();
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
  irsend.sendRaw(rawData, v.size(), 38);
}
