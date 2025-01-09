package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Cell struct {
	Etat                  string //"I" infecté, "S" sain, "G" guéri
	probaInfectionMoyenne float32
	tempsInfectionRestant int //nombre d'itérations où il sera encore infecté
	tempsImmuniteRestant  int
	posi                  int
	posj                  int
}

func initcellule(cellule *Cell, probaInfectionMoyenne float32, tempsInfectionMoyen int, i int, j int, proba float32, wg *sync.WaitGroup) {
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
	wg.Done()
}

func main() {

	var size int
	var proba float32 = 0.2
	var probaInfectionMoyenne float32
	var nbIterations int
	var tempsInfectionMoyen int
	var wg sync.WaitGroup
	fmt.Print("entrez la taille de la grille ")
	fmt.Scanln(&size)
	fmt.Print("Entrez le nombre d'itérations ")
	fmt.Scanln(&nbIterations)
	fmt.Print("Choisissez le temps d'infection moyen ")
	fmt.Scanln(&tempsInfectionMoyen)

	rand.Seed(time.Now().UnixNano())

	matrice := make([][]Cell, size)
	for i := range matrice {
		matrice[i] = make([]Cell, size)

	}
	wg.Add(size * size)

	for i := range matrice {
		for j := range matrice[i] {
			go initcellule(&matrice[i][j], probaInfectionMoyenne, tempsInfectionMoyen, i, j, proba, &wg)

		}
	}
	wg.Wait()
	//programme principal
	for i := 0; i < nbIterations; i++ {
		wg.Add(size * size)

		for j := range matrice {
			for k := range matrice[j] {
				go evolveCell(&matrice[j][k], matrice, &wg)
			}
		}
		wg.Wait()
	}

	fmt.Println("Hello, World!")
}
