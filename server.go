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
	ln, err := net.Listen("tcp", portString) // Écouter sur le port
	if err != nil {
		fmt.Println("Erreur lors de l'écoute sur le port :", err)
		return
	}

	defer ln.Close() // Fermeture de la connexion
	message := fmt.Sprintf("Serveur en écoute sur le port %s", portString)
	fmt.Println(message)

	// Boucle infinie pour accepter toutes les connexions entrantes
	for {
		// Acceptation d'une nouvelle connexion sur ce port
		conn, errconn := ln.Accept()
		if errconn != nil {
			fmt.Println("Erreur de tentative de connexion :", errconn)
			continue
		}
		// Gestion de la connexion dans une goroutine
		go gestionConnexion(conn)
	}
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
		fmt.Println("Erreur lors de la lecture :", err)
		return ""
	}

	// Nettoie la réponse
	reponse = strings.TrimSpace(reponse)
	return reponse
}

func gestionConnexion(conn net.Conn) {
	defer conn.Close() // Assure que la connexion est fermée à la fin de la fonction

	reader := bufio.NewReader(conn)

	// Demande des paramètres au client
	size := demanderAuClient(reader, conn, "Veuillez entrer la taille de la grille :")
	nombreIterations := demanderAuClient(reader, conn, "Veuillez entrer le nombre d'itérations :")
	tempsInfection := demanderAuClient(reader, conn, "Veuillez entrer le temps d'infection moyen :")
	afficherEtat := demanderAuClient(reader, conn, "Voulez-vous afficher l'état de l'automate dans la console à chaque itération ? (oui/non)")
	size_int, mist := strconv.Atoi(size)
	if mist != nil {
		fmt.Println("Erreur de conversion:", mist)
		return
	}
	// Confirmation de réception des paramètres
	fmt.Printf("Paramètres reçus : Taille de la grille = %s, Nombre d'itérations = %s, Temps d'infection = %s, Afficher état = %s\n",
		size, nombreIterations, tempsInfection, afficherEtat)

	// Initialisation de la grille ou autre traitement
	_, err := io.WriteString(conn, "Initialisation de la grille...\n")
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de l'initialisation :", err)
		//affichage de la matrice initiale
		matrice := make([][]Cell, size_int)
		for i := range matrice {
			matrice[i] = make([]Cell, size_int)
		}
		print(MatrixtoString(matrice))
	}
}

func sendMatrix(matrice *matrice, conn net.Conn) {
	_, err := io.WriteString(conn, matrice+"\n")
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la matrice :", err)
		return
	}
}
