
#include <ESP8266httpUpdate.h>

extern bool isInUpgrading;

void update_started() {
  Serial.println("update started");
}

void update_finished() {
  Serial.println("update finished");
}

void update_progress(int cur, int total) {
  Serial.printf("update process at %d of %d bytes...\n", cur, total);
}

void update_error(int err) {
  Serial.printf("update fatal error code %d\n", err);
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
        Serial.printf("update Error (%d): %s\n", ESPhttpUpdate.getLastError(), ESPhttpUpdate.getLastErrorString().c_str());
        r = false;
        break;

      case HTTP_UPDATE_NO_UPDATES:
        Serial.println("no update");
        r = true;
        break;

      case HTTP_UPDATE_OK:
        Serial.println("update ok");
        r = true;
        break;
    }
    isInUpgrading = false;
    return r;
}
