

#include <WiFiNINA.h>  
#include "wifi_connection.h"


WiFiSSLClient wifi;
const char* webSocketServerAddress = "localhost";
const int webSocketServerPort = 8080; // Port du serveur WebSocket
WiFiClient client;

void setup() {
  Serial.begin(9600);
  if (!Serial) delay(3000);
  wifi_connection();

  Serial.println("Connecting to WebSocket server...");
  if (client.connect(webSocketServerAddress, webSocketServerPort)) {
    Serial.println("Connected to WebSocket server");

    // Envoyer une demande de poignée de main (handshake) WebSocket
    client.println("GET /ws HTTP/1.1");
    client.println("Host: localhost");
    client.println("Upgrade: websocket");
    client.println("Connection: Upgrade");
   
    client.println("Sec-WebSocket-Version: 13");
    client.println();
  } else {
    Serial.println("Connection failed");
  }
}

void loop() {
  // Lire et traiter les données du serveur WebSocket
}