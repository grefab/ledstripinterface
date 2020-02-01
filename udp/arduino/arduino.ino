#include <FastLED.h>
#define NUM_LEDS 216 // 144 LEDs/m, 1.5m
#define DATA_PIN 5
CRGB leds[NUM_LEDS];

#include <EtherCard.h>
#include <IPAddress.h>

#define STATIC 0  // set to 1 to disable DHCP (adjust myip/gwip values below)

static byte mymac[] = { 0x70, 0x69, 0x69, 0x2D, 0x30, 0x31 };
#if STATIC
static byte myip[] = { 10, 42, 0, 201 };
static byte gwip[] = { 10, 42, 0, 1 };
#endif
byte Ethernet::buffer[NUM_LEDS * 3 + 43]; // TCP/IP send and receive buffer (frame size), UDP adds 43 bytes of header before data

void setupLeds() {
  LEDS.addLeds<WS2812, DATA_PIN, GRB>(leds, NUM_LEDS);
  LEDS.setBrightness(32);
  initLedsBlack();
  initLedsYellow();
}

void ledWelcome() {
  initLedsWhite();
  delay(100);
  initLedsGreen();
  delay(100);
  initLedsWhite();
  delay(100);
  initLedsGreen();
  delay(100);
  initLedsWhite();
}
void initLedsBlack() {
  for (int i = 0; i < NUM_LEDS; i++) {
    leds[i] = CRGB::Black;
  }
  FastLED.show();
}
void initLedsWhite() {
  for (int i = 0; i < NUM_LEDS; i++) {
    leds[i] = CRGB::White;
    leds[i].fadeLightBy(220);
  }
  FastLED.show();
}
void initLedsGreen() {
  for (int i = 0; i < NUM_LEDS; i++) {
    leds[i] = CRGB::Green;
  }
  FastLED.show();
}
void initLedsYellow() {
  for (int i = 0; i < NUM_LEDS; i++) {
    leds[i] = CRGB::Yellow;
    leds[i].fadeLightBy(220);
  }
  FastLED.show();
}

void udpSerialPrint(uint16_t dest_port,
                    uint8_t src_ip[IP_LEN],
                    uint16_t src_port,
                    const char *data,
                    uint16_t len) {
  for (int i = 0; i < NUM_LEDS; ++i) {
    leds[i] = CRGB(data[i * 3 + 0], data[i * 3 + 1], data[i * 3 + 2]);
  }
  FastLED.show();
}

void setup() {
  setupLeds();
  Serial.begin(57600);
  Serial.println("setting up ethernet");
  if (!ether.begin(sizeof Ethernet::buffer, mymac, SS)) {
    Serial.println("Failed to access Ethernet controller");
  }
#if STATIC
  Serial.print("setting up static ip");
  ether.staticSetup(myip, gwip);
#else
  Serial.print("setting up dhcp");
  if (!ether.dhcpSetup()) {
    Serial.println("DHCP failed");
  }
#endif
  ether.printIp("IP:  ", ether.myip);
  ether.printIp("GW:  ", ether.gwip);
  ether.printIp("DNS: ", ether.dnsip);

  ether.udpServerListenOnPort(&udpSerialPrint, 1337);
  ledWelcome();
}

void loop() {
  ether.packetLoop(ether.packetReceive());
}
