package main

import (
	"image/color"
	"log"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	// Créer une matrice 10x10 avec des valeurs aléatoires
	rows, cols := 10, 10
	matrix := make([][]float64, rows)
	for i := range matrix {
		matrix[i] = make([]float64, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64() // Valeurs aléatoires entre 0 et 1
		}
	}

	// Créer un graphique
	p := plot.New()

	// Créer un tableau de points pour la matrice
	pts := make(plotter.XYs, rows*cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			pts[i*cols+j].X = float64(j)
			pts[i*cols+j].Y = float64(i)
		}
	}

	// Créer un scatter plot
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}

	// Ajouter des couleurs basées sur les valeurs de la matrice
	scatter.GlyphStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Couleur rouge
	scatter.GlyphStyle.Radius = vg.Points(5)

	// Ajouter le scatter plot au graphique
	p.Add(scatter)

	// Définir les limites des axes
	p.X.Min = 0
	p.X.Max = float64(cols)
	p.Y.Min = 0
	p.Y.Max = float64(rows)

	// Sauvegarder le graphique dans un fichier PNG
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "matrix.png"); err != nil {
		log.Fatal(err)
	}
}
