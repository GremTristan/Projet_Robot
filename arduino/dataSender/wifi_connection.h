

const char *WIFI_SSID = "Bbox-AA2146FC";
const char *WIFI_PASSWORD = "6gkV1wET93wp32Tgep";

void wifi_connection(){

 
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  Serial.println("Connecting to Wi-Fi");
  Serial.println(WIFI_SSID);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("Connected to Wi-Fi");
  Serial.print("Adresse IP de l'Arduino : ");
  Serial.println(WiFi.localIP());

  
  }