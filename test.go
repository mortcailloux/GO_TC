package main

import (
	"fmt"
	"math/rand"
	"time"
)

// initcellule séquentielle
func initcellule2(cellule *Cell, probaInfectionMoyenne float32, tempsInfectionMoyen int, i int, j int, proba float32) {
	cellule.posi = i
	cellule.posj = j

	aleatoire := rand.Float32() //genere un nombre aleatoire entre 0 et 1 (loi uniforme)
	if aleatoire > proba {
		cellule.Etat = "S"
		cellule.probaInfectionMoyenne = probaInfectionMoyenne
		cellule.tempsInfectionRestant = -1
		cellule.tempsImmuniteRestant = -1

	} else {
		cellule.Etat = "I"
		cellule.probaInfectionMoyenne = probaInfectionMoyenne
		cellule.tempsInfectionRestant = tempsInfectionMoyen + rand.Intn(11) - 5 // temps moyen + nombre aléatoire entre -5 et 5
		cellule.tempsImmuniteRestant = -1

	}
}

// eveolvecell séquentielle
func evolveCell2(cell *Cell, grid [][]Cell, changement *bool) {
	rand.Seed(time.Now().UnixNano())

	// Find neighboring cells
	neighbors := findNeighbors1(cell, grid)

	// For each possible state
	switch cell.Etat {
	case "S": // Susceptible
		infectedNeighbors := countInfectedNeighbors(neighbors)

		// Calculate infection probability
		// infectionProbability = base infection chance * number of infected neighbors / total neighbors
		proba := cell.probaInfectionMoyenne * float32(infectedNeighbors) / float32(len(neighbors))

		// Infection only happens if the random number is less than the infection probability
		if rand.Float32() < proba {

			*changement = true
			cell.Etat = "I"
			cell.tempsInfectionRestant = rand.Intn(5) + 1 // Random infection duration (e.g., 1-5 iterations)
		}
	case "I": // Infected
		cell.tempsInfectionRestant--
		if cell.tempsInfectionRestant <= 0 {

			*changement = true
			cell.Etat = "G"
			cell.tempsImmuniteRestant = rand.Intn(5) + 1 // Random immunity duration (e.g., 1-5 iterations)
		}
	case "G": // Recovered (immune)
		cell.tempsImmuniteRestant--

		if cell.tempsImmuniteRestant <= 0 {

			*changement = true
			cell.Etat = "S" // Becomes susceptible again
		}
	}
	//utile pour arrêter le programme plus tôt s'il n'y a plus de changements même si les itérations n'ont pas fini
}

// test les performances de l'execution séquentielle
func performances(nbIterations int, size int, tempsInfectionMoyen int, probaInfectionMoyenne float32, proba float32) {
	fmt.Printf("performances")
	start2 := time.Now()
	matrice := make([][]Cell, size)
	for i := range matrice {
		matrice[i] = make([]Cell, size)

	}

	for i := range matrice {
		for j := range matrice[i] {
			initcellule2(&matrice[i][j], probaInfectionMoyenne, tempsInfectionMoyen, i, j, proba)

		}
	}
	fmt.Print("Execution du programme principal\n")
	changement := false
	//programme principal
	var iter int

	for i := 0; i < nbIterations; i++ {

		for j := range matrice {
			for k := range matrice[j] {
				evolveCell2(&matrice[j][k], matrice, &changement)
			}
		}

		iter = i
	}
	if iter < nbIterations-1 {
		fmt.Printf("Le programme a trouvé un état stable et s'est arrêté à la %d itération", iter)

	}
	visualizeMatrix(matrice, "fin.png")
	fmt.Printf("\ntemps d'execution séquentielle: %v", time.Since(start2))
	//programme qui teste les performances
}
