package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Cell struct {
	Etat                  string
	probaInfectionMoyenne float32
	tempsInfectionRestant int
	tempsImmuniteRestant  int
	posi                  int
	posj                  int
}

// Copie profonde d'une cellule (pas seulement le pointeur)
func (c *Cell) clone() Cell {
	return Cell{
		Etat:                  c.Etat,
		probaInfectionMoyenne: c.probaInfectionMoyenne,
		tempsInfectionRestant: c.tempsInfectionRestant,
		tempsImmuniteRestant:  c.tempsImmuniteRestant,
		posi:                  c.posi,
		posj:                  c.posj,
	}
}

// initialise la cellule avec des probas qu'elle soit infectée, le temps qu'elle restera infectée (une immunité variable)
func initcellule(cellule *Cell, probaInfectionMoyenne float32, tempsInfectionMoyen int, i int, j int, proba float32) {
	cellule.posi = i
	cellule.posj = j

	if rand.Float32() > proba {
		cellule.Etat = "S"
		cellule.probaInfectionMoyenne = probaInfectionMoyenne
		cellule.tempsInfectionRestant = -1
		cellule.tempsImmuniteRestant = -1
	} else {
		cellule.Etat = "I"
		cellule.probaInfectionMoyenne = probaInfectionMoyenne
		cellule.tempsInfectionRestant = tempsInfectionMoyen + rand.Intn(11) - 5
		cellule.tempsImmuniteRestant = -1
	}
}

// échange la grille de lecture et d'écriture pour garder des mises à jour cohérentes
func swapGrids(current, next *[][]Cell) {
	*current, *next = *next, *current
}

func main() {
	var size int
	var proba float32 = 0.2
	var probaInfectionMoyenne float32
	var nbIterations, tempsInfectionMoyen int
	var display bool
	var temp, fin string

	fmt.Print("entrez la taille de la grille ")
	fmt.Scanln(&size)
	fmt.Print("Entrez le nombre d'itérations ")
	fmt.Scanln(&nbIterations)
	fmt.Print("Choisissez le temps d'infection moyen ")
	fmt.Scanln(&tempsInfectionMoyen)
	fmt.Print("Voulez-vous afficher l'état de l'automate dans la console à chaque itération ? (oui/non)")
	fmt.Scanln(&temp)
	display = temp == "oui" || temp == "Oui" || temp == "OUI"

	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	// Création des deux grilles pour éviter d'écrire dans la même grile que celle où l'on lit
	currentGrid := make([][]Cell, size)
	nextGrid := make([][]Cell, size)
	for i := range currentGrid {
		currentGrid[i] = make([]Cell, size)
		nextGrid[i] = make([]Cell, size)
		for j := range currentGrid[i] {
			initcellule(&currentGrid[i][j], probaInfectionMoyenne, tempsInfectionMoyen, i, j, proba)
		}
	}

	numWorkers := 8                         //mettre la même valeur que le nombre de coeur
	batchSize := (size * size) / numWorkers //on calcule la taille de la sous grille sur laquelle on va travailler (pas forcément un carré)
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
				end = size * size //si les divisions par numWorkers n'avait pas un reste nul
			}
			go processBatch(currentGrid, nextGrid, start, end, size, &changement, &wg)
		}

		wg.Wait() //on attends que tous les worker polls aient fini avant de recommencer un tour

		// Échange des grilles
		swapGrids(&currentGrid, &nextGrid)

		if display {
			displayMatrix(currentGrid)
		}

	}

	fmt.Printf("\nTemps d'exécution avec goroutines sur plusieurs cases: %v\n", time.Since(start))
	performances(nbIterations, size, tempsInfectionMoyen, probaInfectionMoyenne, proba)

	fmt.Print("Appuyez sur n'importe quelle touche pour quitter le programme")

	fmt.Scanln(&fin)
}
