package main

import (
	"math/rand"
	"sync"
	"time"
)

// Cell represents a single individual in the grid

// evolveCell evolves the state of a cell based on its neighbors
func evolveCell(cell *Cell, grid [][]Cell, wg *sync.WaitGroup, syncmodification *sync.WaitGroup, changement *bool) {
	rand.Seed(time.Now().UnixNano())

	// Find neighboring cells
	neighbors := findNeighbors(cell, grid)

	// For each possible state
	switch cell.Etat {
	case "S": // Susceptible
		infectedNeighbors := countInfectedNeighbors(neighbors)

		// Calculate infection probability
		// infectionProbability = base infection chance * number of infected neighbors / total neighbors
		proba := cell.probaInfectionMoyenne * float32(infectedNeighbors) / float32(len(neighbors))

		// Infection only happens if the random number is less than the infection probability
		syncmodification.Done()
		syncmodification.Wait()
		if rand.Float32() < proba {

			*changement = true
			cell.Etat = "I"
			cell.tempsInfectionRestant = rand.Intn(5) + 1 // Random infection duration (e.g., 1-5 iterations)
		}
	case "I": // Infected
		cell.tempsInfectionRestant--
		syncmodification.Done()
		syncmodification.Wait()
		if cell.tempsInfectionRestant <= 0 {

			*changement = true
			cell.Etat = "G"
			cell.tempsImmuniteRestant = rand.Intn(5) + 1 // Random immunity duration (e.g., 1-5 iterations)
		}
	case "G": // Recovered (immune)
		syncmodification.Done()
		syncmodification.Wait()
		cell.tempsImmuniteRestant--

		if cell.tempsImmuniteRestant <= 0 {

			*changement = true
			cell.Etat = "S" // Becomes susceptible again
		}
	}
	wg.Done()
	//utile pour arrêter le programme plus tôt s'il n'y a plus de changements même si les itérations n'ont pas fini
}

// findNeighbors finds the neighbors of a cell in the grid
func findNeighbors(cell *Cell, grid [][]Cell) []Cell {
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
