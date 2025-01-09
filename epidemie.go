package main

import (
	"image/color"
	"log"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	Sain    = 0 // 0 pour sain
	Infecte = 1 // 1 pour infecté
	Gueri   = 2 // 2 pour guéri
)

type ColoredScatter struct {
	Xs, Ys []float64
	Colors []color.Color
}

func (cs *ColoredScatter) Len() int {
	return len(cs.Xs)
}

func (cs *ColoredScatter) XY(i int) (float64, float64) {
	return cs.Xs[i], cs.Ys[i]
}

func main() {
	// Dimensions de la matrice
	rows, cols := 10, 10
	matrix := make([][]int, rows)

	// Initialiser la matrice avec des personnes saines
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = Sain // Tout le monde est sain au départ
		}
	}

	// Simuler quelques infections
	for i := 0; i < 5; i++ { // 5 personnes infectées
		x := rand.Intn(rows)
		y := rand.Intn(cols)
		matrix[x][y] = Infecte
	}

	// Simuler quelques guérisons
	for i := 0; i < 3; i++ { // 3 personnes guéries
		x := rand.Intn(rows)
		y := rand.Intn(cols)
		if matrix[x][y] == Infecte {
			matrix[x][y] = Gueri
		}
	}

	// Créer un graphique
	p := plot.New()

	// Créer un tableau de points pour la matrice
	pts := make([]float64, rows*cols*2)      // 2 valeurs (X, Y) pour chaque point
	colors := make([]color.Color, rows*cols) // Tableau pour les couleurs
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			index := i*cols + j
			pts[index*2] = float64(j)   // X
			pts[index*2+1] = float64(i) // Y

			// Déterminer la couleur en fonction de l'état de la personne
			switch matrix[i][j] {
			case Sain:
				colors[index] = color.RGBA{G: 255, A: 255} // Vert
			case Infecte:
				colors[index] = color.RGBA{R: 255, A: 255} // Rouge
			case Gueri:
				colors[index] = color.RGBA{B: 255, A: 255} // Bleu
			}
		}
	}

	// Créer un scatter plot personnalisé
	scatter := &ColoredScatter{
		Xs:     pts[0 : len(pts)/2],
		Ys:     pts[len(pts)/2:],
		Colors: colors,
	}

	// Créer un scatter plot avec un style personnalisé
	scatterPlot, err := plotter.NewScatter(scatter)
	if err != nil {
		log.Fatal(err)
	}

	// Ajouter un style aux points
	for i := 0; i < scatter.Len(); i++ {
		scatterPlot.GlyphStyle.Color = scatter.Colors[i]
		scatterPlot.GlyphStyle.Radius = vg.Points(10) // Taille des points
	}

	// Ajouter le scatter plot au graphique
	p.Add(scatterPlot)

	// Définir les limites des axes
	p.X.Min = 0
	p.X.Max = float64(cols)
	p.Y.Min = 0
	p.Y.Max = float64(rows)

	// Sauvegarder le graphique dans un fichier PNG
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "epidemie.png"); err != nil {
		log.Fatal(err)
	}
}
