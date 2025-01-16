package main

import (
	"math/rand"
	"sync"
)

// Cell represents a single individual in the grid

// Calcule le prochain état d'une cellule sans modifier la cellule originale, renvoie une nouvelle Cellule correspondant à la cellule modifiée
func evolveCell(currentCell Cell, neighbors []Cell) Cell {

	newCell := currentCell.clone()

	switch currentCell.Etat {
	case "S":
		infectedCount := 0
		for _, n := range neighbors {
			if n.Etat == "I" {
				infectedCount++
			}
		}
		proba := currentCell.probaInfectionMoyenne * float32(infectedCount) / float32(len(neighbors))
		if rand.Float32() < proba {
			newCell.Etat = "I"
			newCell.tempsInfectionRestant = rand.Intn(5) + 1
		}
	case "I":
		newCell.tempsInfectionRestant--
		if newCell.tempsInfectionRestant <= 0 {
			newCell.Etat = "G"
			newCell.tempsImmuniteRestant = rand.Intn(5) + 1
		}
	case "G":
		newCell.tempsImmuniteRestant--
		if newCell.tempsImmuniteRestant <= 0 {
			newCell.Etat = "S"
		}
	}
	return newCell
}

func findNeighbors(cell *Cell, grid [][]Cell, size int) []Cell {
	directions := []struct{ di, dj int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
	}

	neighbors := make([]Cell, 0, 8)
	for _, d := range directions {
		ni, nj := cell.posi+d.di, cell.posj+d.dj
		if ni >= 0 && ni < size && nj >= 0 && nj < size {
			neighbors = append(neighbors, grid[ni][nj])
		}
	}
	return neighbors
}

func processBatch(currentGrid, nextGrid [][]Cell, start, end, size int, changement *bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := start; i < end; i++ {
		row := i / size //on va parcourir notre nombre de cellules du tableau de gauche à droite puis de haut en bas
		//cellule 0 par exemple avec, cela donne ligne 0 et col 0, cellule 1, ligne0, col 1 (puisque / est le quotient et % le reste)
		//c'est un peu des maths mais ça fonctionne
		col := i % size

		// Lecture des voisins depuis la grille de lecture (on écrit dans la 2e pour éviter les conflits)
		neighbors := findNeighbors(&currentGrid[row][col], currentGrid, size)

		// Calcul du nouvel état
		newCell := evolveCell(currentGrid[row][col], neighbors)

		// Si l'état a changé, on le note
		if newCell.Etat != currentGrid[row][col].Etat {
			*changement = true
		}
		//la variable changement ne fonctionnait pas puisqu'elle est accédée par de nombreux processus, le programme ne fonctonnait pas correctement,
		//à voir si ça vaut le coup de patch
		// Écriture dans la nouvelle grille (grille d'écriture)
		nextGrid[row][col] = newCell
	}
}
func findNeighbors1(cell *Cell, grid [][]Cell) []Cell {
	directions := []struct{ di, dj int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // Cardinal directions (up, down, left, right)
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, // Diagonal directions
	}

	var neighbors []Cell
	rows := len(grid)
	cols := len(grid[0])

	for _, d := range directions {
		ni, nj := cell.posi+d.di, cell.posj+d.dj
		if ni >= 0 && ni < rows && nj >= 0 && nj < cols {
			neighbors = append(neighbors, grid[ni][nj])
		}
	}

	return neighbors
}

// countInfectedNeighbors counts the number of infected neighbors
func countInfectedNeighbors(neighbors []Cell) int {
	count := 0
	for _, neighbor := range neighbors {
		if neighbor.Etat == "I" {
			count++
		}
	}
	return count
}
