package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func client(portString string) {
	var fin string
	var buffer strings.Builder
	// Connect to the server
	conn, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to the server.")
	// Reader for server responses
	serverReader := bufio.NewReader(conn)
	// Reader for user input
	userInputReader := bufio.NewReader(os.Stdin)
	//l'échange se fait en 2 phases, une phase questions/réponses entre le client et le serveur et une phase où le serveur
	//fournit des données au client
	// Phase 1 : Répondre aux questions
	for {
		message, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Erreur de lecture :", err)
			return
		}

		fmt.Print("Serveur : " + message)

		// Détection du début de la Phase 2
		if strings.TrimSpace(message) == "Paramètres reçus, début de l'envoi des données." { //fin de la phase 1
			break
		}

		// Répondre à la question
		fmt.Print(">> ")
		response, err := userInputReader.ReadString('\n')
		if err != nil {
			fmt.Println("Erreur lors de la lecture de l'entrée utilisateur :", err)
			return
		}

		_, err = io.WriteString(conn, response)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi au serveur :", err)
			return
		}
	}

	// Phase 2 : Recevoir des données en continu
	for {
		// on lit une ligne reçue
		data, err := serverReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Le serveur a terminé l'envoi des données.")
				break
			}
			fmt.Println("Erreur de lecture des données :", err)
			continue
		}
		// Détection de la fin de l'envoi
		if strings.TrimSpace(data) == "FIN_DATA" {
			fmt.Println("Fin de l'envoi des données.") //fin du programme
			break
		}

		// Ajouter la ligne au buffer
		buffer.WriteString(data)
		// Vérifier si le buffer contient "G:Guéri" (fin du bloc à afficher)
		if strings.Contains(buffer.String(), "G: Guéri") {
			// Afficher tout le contenu accumulé jusqu'à présent
			fmt.Println("Bloc complet reçu :")
			fmt.Print(buffer.String())

			// Réinitialiser le buffer pour recevoir les nouveaux blocs
			buffer.Reset()
		}

	}

	// Ajouter une pause avant de quitter le programme
	fmt.Print("Appuyez sur Entrée pour terminer")
	fmt.Scanln(&fin)
}
