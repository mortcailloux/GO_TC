package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
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
	nbIterations, err := strconv.Atoi(nombreIterations)
	if err != nil {
		fmt.Println("Erreur de conversion du nombre d'itérations :", err)
		_, _ = io.WriteString(conn, "Erreur de conversion de la taille de la grille.\n")
		return
	}
	tempsInfectionMoyen, err := strconv.Atoi(nombreIterations)
	if err != nil {
		fmt.Println("Erreur de conversion du temps d'infection :", err)
		_, _ = io.WriteString(conn, "Erreur de conversion de la taille de la grille.\n")
		return
	}
	// Confirmation des réponses reçues
	fmt.Printf("Paramètres reçus : Taille = %s, Iterations = %s, TempsInfection = %s\n", size, nombreIterations, tempsInfection)
	_, _ = io.WriteString(conn, "Paramètres reçus, début de l'envoi des données.\n")

	// Phase 2 : Envoi continu de données

	var proba float32 = 0.2
	var probaInfectionMoyenne float32

	var display bool

	display = true

	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	// Création des deux grilles pour éviter d'écrire dans la même grile que celle où l'on lit
	currentGrid := make([][]Cell, sizeInt)
	nextGrid := make([][]Cell, sizeInt)
	for i := range currentGrid {
		currentGrid[i] = make([]Cell, sizeInt)
		nextGrid[i] = make([]Cell, sizeInt)
		for j := range currentGrid[i] {
			initcellule(&currentGrid[i][j], probaInfectionMoyenne, tempsInfectionMoyen, i, j, proba)
		}
	}

	numWorkers := 8                               //mettre la même valeur que le nombre de coeur
	batchSize := (sizeInt * sizeInt) / numWorkers //on calcule la taille de la sous grille sur laquelle on va travailler (pas forcément un carré)
	if batchSize < 1 {
		batchSize = 1
	}

	var wg sync.WaitGroup
	changement := false

	// Boucle principale
	for iter := 0; iter < nbIterations; iter++ {
		changement = false
		wg.Add(numWorkers)

		for w := 0; w < numWorkers; w++ {
			start := w * batchSize   //on calcule la première cellule sur laquelle on va effectuer des modifications
			end := start + batchSize //on calcule la dernière cellule
			if w == numWorkers-1 {
				end = sizeInt * sizeInt //si les divisions par numWorkers n'avait pas un reste nul
			}
			go processBatch(currentGrid, nextGrid, start, end, sizeInt, &changement, &wg)
		}

		wg.Wait() //on attends que tous les worker polls aient fini avant de recommencer un tour

		// Échange des grilles
		swapGrids(&currentGrid, &nextGrid)

		if display {
			sendMatrix(&currentGrid, conn)
		}

	}
	// Signal de fin
	_, _ = io.WriteString(conn, "FIN_DATA\n")

	fmt.Printf("\nTemps d'exécution avec goroutines sur plusieurs cases: %v\n", time.Since(start))
	performances(nbIterations, sizeInt, tempsInfectionMoyen, probaInfectionMoyenne, proba)

}

func sendMatrix(matrice *[][]Cell, conn net.Conn) {
	fmt.Println("Envoie de la matrice ")
	stringMatrix := MatrixtoString(matrice)
	fmt.Print(stringMatrix)
	_, err := io.WriteString(conn, stringMatrix+"\n")
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la matrice :", err)
		return
	}
}
