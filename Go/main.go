package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type cell struct {
	Etat                  string //"I" infecté, "S" sain, "G" guéri
	probaInfectionMoyenne float32
	tempsInfectionRestant int //nombre d'itérations où il sera encore infecté
	tempsImmuniteRestant  int
	posi                  int
	posj                  int
}

func cellule(carre *cell) {

}

func initcellule(cellule *cell, probaInfectionMoyenne float32, tempsInfectionMoyen int, i int, j int, proba float32, wg *sync.WaitGroup) {
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
	var size int = 4
	var proba float32 = 0.2
	var probaInfectionMoyenne float32 = 0.3
	var nbIterations int = 1000
	var tempsInfectionMoyen int = 10
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	matrice := make([][]cell, size)
	for i := range matrice {
		matrice[i] = make([]cell, size)

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
				go cellule(&matrice[j][k])
			}
		}
		wg.Wait()
	}

	fmt.Println("Hello, World!")
}
