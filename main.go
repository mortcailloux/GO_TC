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
	var syncmodification sync.WaitGroup
	var display bool
	var temp string
	var fin string
	var changement bool
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
	fmt.Print("Initialisation de la grille...")

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
	fmt.Print("Execution du programme principal")

	//programme principal
	i := 0
	for i := 0; i < nbIterations && !changement; i++ {
		wg.Add(size * size)
		syncmodification.Add(size * size)
		changement = false

		for j := range matrice {
			for k := range matrice[j] {
				go evolveCell(&matrice[j][k], matrice, &wg, &syncmodification, &changement)
			}
		}
		wg.Wait()
		if display {
			displayMatrix(matrice)
		}
	}
	if i < nbIterations {
		fmt.Printf("Le programme a trouvé un état stable et s'est arrêté à la %d itération", i)

	}
	visualizeMatrix(matrice, "fin.png")
	fmt.Print("Appuyez sur n'importe quelle touche pour quitter le programme")
	fmt.Scanln(&fin)
}
