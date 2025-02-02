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
		fmt.Print("\n Veuillez fermer votre serveur précédent")
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
func demanderAuClient(reader *bufio.Reader, conn net.Conn, demande string) (string, error) {
	// Envoie une demande au client
	_, err := io.WriteString(conn, demande+"\n")
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la demande :", err)
		return "", err
	}

	// Lit la réponse du client
	reponse, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur lors de la lecture :", err)
		return "", err
	}

	// Nettoie la réponse
	reponse = strings.TrimSpace(reponse)
	return reponse, err
}

func ask_int(reader *bufio.Reader, conn net.Conn, demande string) (int, error) {

	size, err2 := demanderAuClient(reader, conn, demande)
	sizeInt, err := strconv.Atoi(size)
	for err != nil {
		if err2 != nil {
			fmt.Print("\nErreur lors de la communication au client, abandon de la communication avec lui\n")
			break
		}
		fmt.Println("Erreur de conversion de la taille :", err)
		size, err2 = demanderAuClient(reader, conn, "veuillez entrer un entier")
		sizeInt, err = strconv.Atoi(size)

	}
	return sizeInt, err2
}
func ask_string(reader *bufio.Reader, conn net.Conn, demande string) (string, error) {

	rep, err := demanderAuClient(reader, conn, demande)
	for rep != "oui" && rep != "non" {
		if err != nil {
			fmt.Print("Erreur lors de la communication au client, abandon de la communication avec lui")
			break
		}
		rep, err = demanderAuClient(reader, conn, "veuillez entrer une réponse valide (oui ou non)")

	}
	return rep, err
}

func gestionConnexion(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Phase 1 : Questions initiales
	sizeInt, err := ask_int(reader, conn, "Veuillez entrer la taille de la grille :")
	if err != nil {
		return

	}
	nbIterations, err := ask_int(reader, conn, "Veuillez entrer le nombre d'itérations :")
	if err != nil {
		return
	}
	tempsInfectionMoyen, err := ask_int(reader, conn, "Veuillez entrer le temps d'infection moyen : (nombre d'itérations)")
	if err != nil {
		return
	}
	reponse, err := ask_string(reader, conn, "Voulez-vous enregistrer une image de la grille pour chaque itération ? (oui/non)")
	if err != nil {
		return
	}
	// Conversion de la taille de la grille en entier

	// Confirmation des réponses reçues
	fmt.Printf("Paramètres reçus : Taille = %d, Iterations = %d, TempsInfection = %d\n", sizeInt, nbIterations, tempsInfectionMoyen)
	_, _ = io.WriteString(conn, "Paramètres reçus, début de l'envoi des données.\n")

	// Phase 2 : Envoi continu de données

	var proba float32 = 0.2
	var probaInfectionMoyenne float32

	var display bool
	output := reponse == "oui"
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
		if output {
			nom := fmt.Sprintf("Etape %d.png", iter)
			visualizeMatrix(&currentGrid, nom)

		}

	}
	// Signal de fin
	_, _ = io.WriteString(conn, "FIN_DATA\n")
	fmt.Printf("\nTemps d'exécution avec goroutines sur plusieurs cases: %v\n", time.Since(start))
	if !output {
		visualizeMatrix(&currentGrid, "fin.png")
	}

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
