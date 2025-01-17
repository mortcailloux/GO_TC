package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func server(portString string) {
	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Println("Erreur lors de l'écoute sur le port :", err)
		return
	}
	defer ln.Close()

	fmt.Printf("Serveur en écoute sur %s\n", portString)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erreur de tentative de connexion :", err)
			continue
		}
		go gestionConnexion(conn)
	}
}

func gestionConnexion(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Phase 1 : Questions initiales
	size := demanderAuClient(reader, conn, "Veuillez entrer la taille de la grille :")
	nombreIterations := demanderAuClient(reader, conn, "Veuillez entrer le nombre d'itérations :")
	tempsInfection := demanderAuClient(reader, conn, "Veuillez entrer le temps d'infection moyen :")

	// Conversion de la taille de la grille en entier
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		fmt.Println("Erreur de conversion de la taille :", err)
		_, _ = io.WriteString(conn, "Erreur de conversion de la taille de la grille.\n")
		return
	}

	// Confirmation des réponses reçues
	fmt.Printf("Paramètres reçus : Taille = %s, Iterations = %s, TempsInfection = %s\n", size, nombreIterations, tempsInfection)
	_, _ = io.WriteString(conn, "Paramètres reçus, début de l'envoi des données.\n")

	// Phase 2 : Envoi continu de données
	for i := 0; i < 10; i++ { // Simulation de 10 envois
		data := fmt.Sprintf("Données %d : Simulation de calcul avec grille %dx%d...\n", i+1, sizeInt, sizeInt)
		//ici on va faire l'automate cellulaire
		//il faut faire une fonction autre que main qui va prendre les paramètres pour construire la matrice et envoyer le string qui sera affiché
		//côté client
	}

	// Signal de fin
	_, _ = io.WriteString(conn, "FIN_DATA\n")
}

func demanderAuClient(reader *bufio.Reader, conn net.Conn, demande string) string {
	// Envoie une demande au client
	_, err := io.WriteString(conn, demande+"\n")
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la demande :", err)
		return ""
	}

	// Lit la réponse du client
	reponse, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse :", err)
		return ""
	}

	// Nettoie la réponse
	reponse = strings.TrimSpace(reponse)
	return reponse
}
