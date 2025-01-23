// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"net"
// 	"strings"
// )

// // Task structure to hold client parameters
// type Task struct {
// 	conn             net.Conn
// 	tailleGrille     string
// 	nombreIterations string
// 	tempsInfection   string
// 	afficherEtat     string
// }

// // Worker pool size
// const WorkerPoolSize = 5

// // Main server function
// func server(portString string) {
// 	ln, err := net.Listen("tcp", portString) // Listen on the port
// 	if err != nil {
// 		fmt.Println("Erreur lors de l'écoute sur le port :", err)
// 		return
// 	}
// 	defer ln.Close() // Close the listener when done

// 	fmt.Printf("Serveur en écoute sur le port %s\n", portString)

// 	// Create task channel for the worker pool
// 	taskChan := make(chan Task)

// 	// Start the worker pool
// 	for i := 0; i < WorkerPoolSize; i++ {
// 		go worker(i, taskChan)
// 	}

// 	// Infinite loop to accept incoming connections
// 	for {
// 		conn, err := ln.Accept() // Accept a new connection
// 		if err != nil {
// 			fmt.Println("Erreur de tentative de connexion :", err)
// 			continue
// 		}

// 		// Handle the connection in a goroutine
// 		go gestionConnection(conn, taskChan)
// 	}
// }

// // Handle individual client connection
// func gestionConnection(conn net.Conn, taskChan chan Task) {
// 	defer conn.Close() // Ensure the connection is closed

// 	reader := bufio.NewReader(conn)

// 	// Get client parameters
// 	tailleGrille := demanderAuClient(reader, conn, "Veuillez entrer la taille de la grille :")
// 	nombreIterations := demanderAuClient(reader, conn, "Veuillez entrer le nombre d'itérations :")
// 	tempsInfection := demanderAuClient(reader, conn, "Veuillez entrer le temps d'infection moyen :")
// 	afficherEtat := demanderAuClient(reader, conn, "Voulez-vous afficher l'état de l'automate dans la console à chaque itération ? (oui/non)")

// 	// Create a task with the client's parameters
// 	task := Task{
// 		conn:             conn,
// 		tailleGrille:     tailleGrille,
// 		nombreIterations: nombreIterations,
// 		tempsInfection:   tempsInfection,
// 		afficherEtat:     afficherEtat,
// 	}

// 	// Send the task to the worker pool
// 	taskChan <- task
// }

// // Worker function to process tasks
// func worker(id int, taskChan chan Task) {
// 	for task := range taskChan {
// 		fmt.Printf("[Worker %d] Processing task for client...\n", id)

// 		// Simulate task processing
// 		result := fmt.Sprintf("Grille: %s, Iterations: %s, Temps Infection: %s, Afficher: %s",
// 			task.tailleGrille, task.nombreIterations, task.tempsInfection, task.afficherEtat)

// 		// Send the result back to the client
// 		_, err := io.WriteString(task.conn, "Résultat du traitement : "+result+"\n")
// 		if err != nil {
// 			fmt.Println("Erreur lors de l'envoi du résultat :", err)
// 		}

// 		fmt.Printf("[Worker %d] Task completed for client\n", id)
// 	}
// }

// // Helper function to prompt client for input
// func demanderAuClient(reader *bufio.Reader, conn net.Conn, demande string) string {
// 	_, err := io.WriteString(conn, demande+"\n")
// 	if err != nil {
// 		fmt.Println("Erreur lors de l'envoi de la demande :", err)
// 		return ""
// 	}

// 	reponse, err := reader.ReadString('\n')
// 	if err != nil {
// 		fmt.Println("Erreur lors de la lecture :", err)
// 		return ""
// 	}

// 	return strings.TrimSpace(reponse)
// }