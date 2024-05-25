package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Définition de l'upgrader WebSocket pour upgrader les connexions HTTP vers WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin permet de bypasser les vérifications d'origine pour les WebSockets (à ne pas faire en production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Map pour stocker les clients connectés
var clients = make(map[*websocket.Conn]bool)

// Fonction qui lit les messages des clients WebSocket et les diffuse à tous les autres clients connectés
func reader(conn *websocket.Conn) {
	for {
		// Lire le message envoyé par le client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			// Supprimer le client de la map lorsqu'il se déconnecte
			delete(clients, conn)
			return
		}

		// Afficher le message reçu dans les logs
		log.Println("Message from client:", string(p))

		// Envoyer le message à tous les clients connectés
		for client := range clients {
			err = client.WriteMessage(messageType, p)
			if err != nil {
				log.Println(err)
				// Supprimer le client de la map s'il y a une erreur lors de l'envoi
				delete(clients, client)
				// Fermer la connexion avec le client
				client.Close()
			}
		}
	}
}

// Gestionnaire pour la page d'accueil HTTP
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home HTTP")
}

// Gestionnaire pour les connexions WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrader la connexion HTTP en WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	// Fermez la connexion WebSocket une fois la fonction terminée
	defer conn.Close()

	// Ajouter le client à la map des clients connectés
	log.Println("Client Connect")
	clients[conn] = true

	// Envoyer un message au client une fois connecté
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
		return
	}

	// Appeler la fonction reader pour gérer les messages du client
	reader(conn)
}

// Définir les routes de l'application
func setupRoutes() {
	// Route pour la page d'accueil
	http.HandleFunc("/", homePage)
	// Route pour gérer les connexions WebSocket
	http.HandleFunc("/ws", handleWebSocket)
	// Route pour servir les fichiers statiques (comme index.html)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

// Fonction principale de l'application
func main() {
	fmt.Println("Start playing the game")
	fmt.Println("This is all routes: / and /ws and /static for index.html")
	// Configurer les routes de l'application
	setupRoutes()
	// Démarrer le serveur HTTP sur le port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
