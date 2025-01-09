package main

import (
	"math/rand"
	"time"
)

// Cell represents a single individual in the grid
type Cell struct {
	Etat                  string  // "I" = infected, "S" = susceptible, "G" = recovered
	ProbaInfectionMoyenne float32 // Average infection probability
	TempsInfectionRestant int     // Remaining iterations of infection
	TempsImmuniteRestant  int     // Remaining iterations of immunity
	Posi                  int     // Row position in the grid
	Posj                  int     // Column position in the grid
}

// evolveCell evolves the state of a cell based on its neighbors
// evolveCell evolves the state of a cell based on its neighbors
func evolveCell(cell *Cell, grid [][]Cell) {
	rand.Seed(time.Now().UnixNano())

	// Find neighboring cells
	neighbors := findNeighbors(cell, grid)

	// For each possible state
	switch cell.Etat {
	case "S": // Susceptible
		infectedNeighbors := countInfectedNeighbors(neighbors)

		// Calculate infection probability
		// infectionProbability = base infection chance * number of infected neighbors / total neighbors
		proba := cell.ProbaInfectionMoyenne * float32(infectedNeighbors) / float32(len(neighbors))

		// Infection only happens if the random number is less than the infection probability
		if rand.Float32() < proba {
			cell.Etat = "I"
			cell.TempsInfectionRestant = rand.Intn(5) + 1 // Random infection duration (e.g., 1-5 iterations)
		}
	case "I": // Infected
		cell.TempsInfectionRestant--
		if cell.TempsInfectionRestant <= 0 {
			cell.Etat = "G"
			cell.TempsImmuniteRestant = rand.Intn(5) + 1 // Random immunity duration (e.g., 1-5 iterations)
		}
	case "G": // Recovered (immune)
		cell.TempsImmuniteRestant--
		if cell.TempsImmuniteRestant <= 0 {
			cell.Etat = "S" // Becomes susceptible again
		}
	}
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
		ni, nj := cell.Posi+d.di, cell.Posj+d.dj
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
