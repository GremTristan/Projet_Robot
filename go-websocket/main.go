package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home HTTP")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Supprimer la ligne suivante
	//fmt.Fprintf(w, "Hello WebSocket")

	// Mettre à niveau la connexion HTTP en WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close() // Fermez la connexion WebSocket une fois la fonction terminée

	log.Println("Client Connect")

	// Envoyez un message au client
	err = ws.WriteMessage(websocket.TextMessage, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
		return
	}
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", handleWebSocket)
	// Ajoutez ceci à setupRoutes() pour servir les fichiers statiques.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

}

func main() {
	fmt.Println("Start playing the game")
	fmt.Println("This is all routes : / and /ws and /static for index.html")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var clients = make(map[*websocket.Conn]bool) // Map pour stocker les clients connectés

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Ajouter le client à la map des clients connectés
	clients[conn] = true

	// Boucle pour lire les messages envoyés par le client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(clients, conn) // Supprimer le client de la map lorsqu'il se déconnecte
			return
		}

		// Convertir les données du message en string
		message := string(p)
		log.Println("Message from client:", message)

		// Envoyer le message à tous les clients connectés
		for client := range clients {
			err = client.WriteMessage(messageType, p)
			if err != nil {
				log.Println(err)
				delete(clients, client) // Supprimer le client de la map s'il y a une erreur lors de l'envoi
				client.Close()          // Fermer la connexion avec le client
			}
		}
	}
}
