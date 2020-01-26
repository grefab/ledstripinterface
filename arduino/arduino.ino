#include <FastLED.h>

#define NUM_LEDS 90
#define DATA_PIN 5
#define FRAME_SIZE 3

CRGB leds[NUM_LEDS];

void setup() {
  Serial.begin(57600);
  Serial.print("READY\n");
  LEDS.addLeds<WS2812, DATA_PIN, GRB>(leds, NUM_LEDS);
  LEDS.setBrightness(32);
  initLedsWhite();
  FastLED.show();
  delay(100);
  initLedsGreen();
  FastLED.show();
  delay(100);
  initLedsWhite();
  FastLED.show();
  delay(100);
  initLedsGreen();
  FastLED.show();
  delay(100);
  initLedsWhite();
  FastLED.show();
}

void initLedsWhite() {
  for (int i = 0; i < NUM_LEDS; i++) {
    leds[i] = CRGB::White;
    leds[i].fadeLightBy(220);
  }
}
void initLedsGreen() {
  for (int i = 0; i < NUM_LEDS; i++) {
    leds[i] = CRGB::Green;
  }
}

// a sequence is complete, if we received frameSize 0s.
uint8_t history[FRAME_SIZE] = {255, 255, 255};
int historyPtr = 0;
bool isSequenceComplete(uint8_t b) {
  history[historyPtr++] = b;
  historyPtr %= FRAME_SIZE;
  return history[0] == 0 && history[1] == 0 && history[2] == 0;
}

int idx = 0;
uint8_t data[FRAME_SIZE];
int dataPtr = 0;
void handleByte(uint8_t b) {
  data[dataPtr++] = b;
  if (dataPtr == FRAME_SIZE) {
    // frame complete
    if (idx < NUM_LEDS) {
      leds[idx++ % NUM_LEDS] = CRGB(data[0], data[1], data[2]);
    }
    dataPtr = 0;
  }

  if (isSequenceComplete(b)) {
    FastLED.show();
    idx = 0;
    dataPtr = 0;
    return;
  }
}

void loop() {
  while (Serial.available() > 0) {
    int incomingByte = Serial.read();
    handleByte(incomingByte);
  }
}
