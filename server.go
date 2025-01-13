package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func server(address string, portString string) {
	ln, err := net.Listen("tcp", portString) //écouter sur le port portString
	if err != nil {
		fmt.Println("Erreur lors de l'ecoute sur le port:", err)
		return
	}

	defer ln.Close() //fermeture de la connexion
	message := fmt.Sprintf("serveur en écoute sur le port %d", portString)
	fmt.Println(message)

	// boucle infinie pour accepter toutes les connexions entrantes
	for {
		//acceptation d'une nouvelle connexion sur ce port
		conn, errconn := ln.Accept()
		if errconn != nil {
			fmt.Println("erreur de tentative de connexion", errconn)
			continue
		}
		// gestion de la connection dans une goroutine
		go gestionConnexion(conn)

	}
}

func gestionConnexion(conn net.Conn) {
	/*
		Dans chaque goroutine, le serveur lit et envoie des bytes au client
	*/
	defer conn.Close() // Assure que la connexion est fermée à la fin de la fonction

	// Crée un nouveau lecteur pour lire les données de la connexion
	reader := bufio.NewReader(conn)

	// Lit une ligne de données envoyées par le client
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur lors de la lecture :", err)
		return // Quitte la fonction si une erreur se produit
	}

	// Affiche le message reçu
	fmt.Printf("Données reçues : %s", message)

	// Envoie une réponse au client
	response := fmt.Sprintf("Coucou, vous avez envoyé : %s", message)
	_, err = io.WriteString(conn, response)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la réponse :", err) // Affiche l'erreur si l'envoi échoue
	}
}
