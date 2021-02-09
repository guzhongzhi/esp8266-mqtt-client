
#include <ESP8266httpUpdate.h>

extern bool isInUpgrading;

void update_started() {
  Serial.println("started");
}

void update_finished() {
  Serial.println("finished");
}

void update_progress(int cur, int total) {
  Serial.printf("at %d of %d \n", cur, total);
}

void update_error(int err) {
  Serial.printf("error %d\n", err);
}

bool upgrade(const char* url) {
    isInUpgrading = true;
    ESPhttpUpdate.setLedPin(LED_BUILTIN, LOW);
    // Add optional callback notifiers
    ESPhttpUpdate.onStart(update_started);
    ESPhttpUpdate.onEnd(update_finished);
    ESPhttpUpdate.onProgress(update_progress);
    ESPhttpUpdate.onError(update_error);
    Serial.println("url:");
    Serial.println(url);
    WiFiClient client;
    t_httpUpdate_return ret = ESPhttpUpdate.update(client, url);
    bool r = true;
    switch (ret) {
      case HTTP_UPDATE_FAILED:
        Serial.printf("500 (%d): %s\n", ESPhttpUpdate.getLastError(), ESPhttpUpdate.getLastErrorString().c_str());
        r = false;
        break;

      case HTTP_UPDATE_NO_UPDATES:
        Serial.println("304");
        r = true;
        break;

      case HTTP_UPDATE_OK:
        Serial.println("ok");
        r = true;
        break;
    }
    isInUpgrading = false;
    return r;
}
