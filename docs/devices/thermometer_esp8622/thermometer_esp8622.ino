#include "config.h"
#include <ESP8266WiFi.h>
#include <Wire.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_BME280.h>

// -- program --
ADC_MODE(ADC_VCC);

#ifdef DEBUG_ESP_PORT
#define DEBUG_MSG(...) DEBUG_ESP_PORT.print( __VA_ARGS__ )
#define DEBUG_MSG_F(...) DEBUG_ESP_PORT.printf( __VA_ARGS__ )
#define DEBUG_MSG_LN(...) DEBUG_ESP_PORT.println( __VA_ARGS__ )
#else
#define DEBUG_MSG(...)
#define DEBUG_MSG_F(...)
#define DEBUG_MSG_LN(...)
#endif

unsigned long loopTiming = 0;

Adafruit_BME280 bme; // I2C

void setup() {
  #ifdef DEBUG_ESP_PORT
    Serial.begin(SERIAL_SPEED);
    while (!Serial) {
      yield();
    };
  #endif

  Wire.begin(D2, D1);
  Wire.setClock(100000);

  if (!bme.begin(0x76)) {
    Serial.println("Could not find a valid BME280 sensor, check wiring!");
    while (1);
  }
  DEBUG_MSG_LN("BME280 connected");

  WiFi.mode(WIFI_STA);
  WiFi.persistent(false);

  loopTiming = millis() - LOOP_DELAY;
}

void loop() {
  if (!WiFi.isConnected()) {
    long wifiTimeout = millis();
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
    
    //DEBUG_MSG("WiFi connecting ");
    
    while (WiFi.status() != WL_CONNECTED && (millis() - wifiTimeout < WIFI_CONNECTION_TIMEOUT_MS)) {
      //DEBUG_MSG(".");
      yield();
    }
    //DEBUG_MSG_LN();

    if(WiFi.status() != WL_CONNECTED) {
      //DEBUG_MSG_LN("WiFi connecting failed");
      return;
    }
  
    //DEBUG_MSG_LN();
    //DEBUG_MSG("Connected, IP address: ");
    //DEBUG_MSG_LN(WiFi.localIP());
  
    #ifdef DEBUG_ESP_PORT
      //WiFi.printDiag(Serial);
    #endif
  }

  unsigned long currentMillis = millis();
  if (currentMillis < loopTiming + LOOP_DELAY) {
    return;
  }
  loopTiming = currentMillis;

  //float volts = (float) ESP.getVcc() / 1023.00;
  //DEBUG_MSG("VCC: ");
  //DEBUG_MSG_LN(String(volts, 2));

  DEBUG_MSG("Free heap: ");
  DEBUG_MSG_LN(ESP.getFreeHeap());

  float temperature = bme.readTemperature();
  DEBUG_MSG("Temperature: ");
  DEBUG_MSG_LN(String(temperature, 2));
}
