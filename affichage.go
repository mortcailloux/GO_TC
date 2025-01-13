package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func displayMatrix(matrix [][]Cell) error { //affichage matrice console (avec I/G/S)
	rows := len(matrix)
	if rows == 0 {
		return fmt.Errorf("matrice vide")
	}

	fmt.Println("État actuel de la matrice :")
	for i := range matrix {
		for j := range matrix[i] {
			switch matrix[i][j].Etat {
			case "S":
				fmt.Print("S ")
			case "I":
				fmt.Print("I ")
			case "G":
				fmt.Print("G ")
			default:
				fmt.Print("? ")
			}
		}
		fmt.Println() // Nouvelle ligne pour chaque rangée
	}

	// Afficher la légende
	fmt.Println("\nLégende:")
	fmt.Println("S: Sain")
	fmt.Println("I: Infecté")
	fmt.Println("G: Guéri")
	fmt.Println() //espaces à la fin pour distinguer les matrices quand elles sont affichées à la suite
	fmt.Println()
	fmt.Println()
	return nil
}

func MatrixtoString(matrix [][]Cell) string {
	rows := len(matrix)
	if rows == 0 {
		return "matrice vide"
	}
	retour := ""

	fmt.Println("État actuel de la matrice :")
	for i := range matrix {
		for j := range matrix[i] {
			switch matrix[i][j].Etat {
			case "S":
				retour += "S "
			case "I":
				retour += "I "
			case "G":
				retour += "G "
			default:
				retour += "? "
			}
		}
		retour += "\n" // Nouvelle ligne pour chaque rangée
	}

	// Afficher la légende
	retour += "\nLégende:"
	retour += "S: Sain"
	retour += "I: Infecté"
	retour += "G: Guéri"
	retour += "\n \n \n " //espaces à la fin pour distinguer les matrices quand elles sont affichées à la suite

	return retour

}

func createSquare(x, y float64) plotter.XYs {
	return plotter.XYs{
		{X: x - 0.5, Y: y - 0.5}, // Coin inférieur gauche
		{X: x + 0.5, Y: y - 0.5}, // Coin inférieur droit
		{X: x + 0.5, Y: y + 0.5}, // Coin supérieur droit
		{X: x - 0.5, Y: y + 0.5}, // Coin supérieur gauche
	}
}

func visualizeMatrix(matrix [][]Cell, filename string) error {
	rows := len(matrix)
	if rows == 0 {
		return fmt.Errorf("matrice vide")
	}
	cols := len(matrix[0])

	// Créer un graphique
	p := plot.New()

	// Pour chaque cellule de la matrice
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			// Créer un polygone carré
			square := createSquare(float64(j), float64(i))
			poly, err := plotter.NewPolygon(square)
			if err != nil {
				return err
			}

			// Définir la couleur selon l'état
			switch matrix[i][j].Etat {
			case "S":
				poly.Color = color.RGBA{G: 255, A: 255} // Vert pour Sain
			case "I":
				poly.Color = color.RGBA{R: 255, A: 255} // Rouge pour Infecté
			case "G":
				poly.Color = color.RGBA{B: 255, A: 255} // Bleu pour Guéri
			}

			// Ajouter le polygone au graphique
			p.Add(poly)
		}
	}

	// Configurer le graphique
	p.Title.Text = "État des Cellules"
	p.X.Label.Text = "Position X"
	p.Y.Label.Text = "Position Y"

	// Définir les limites des axes
	p.X.Min = -1
	p.X.Max = float64(cols)
	p.Y.Min = -1
	p.Y.Max = float64(rows)

	// Sauvegarder le graphique
	return p.Save(8*vg.Inch, 8*vg.Inch, filename)
}
func test() {
	// Créer une matrice 10x10 de cellules
	rows, cols := 10, 10
	matrix := make([][]Cell, rows)
	for i := range matrix {
		matrix[i] = make([]Cell, cols)
		for j := range matrix[i] {
			// Attribution aléatoire des états pour la démonstration
			etat := ""
			switch rand.Intn(3) {
			case 0:
				etat = "S" // Sain
			case 1:
				etat = "I" // Infecté
			case 2:
				etat = "G" // Guéri
			}
			matrix[i][j] = Cell{
				Etat: etat,
				posi: i,
				posj: j,
			}
		}
	}

	// Visualiser la matrice
	if err := visualizeMatrix(matrix, "matrix.png"); err != nil {
		log.Fatal(err)
	}
}
