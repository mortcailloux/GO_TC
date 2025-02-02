repo pour GO en tc

### description du projet
Le projet consiste en un automate cellulaire pour simuler la propagation d'une épidémie dans une population.
### fonctionnement
On laisse l'utilisateur décider de la durée de la simulation (le nombre d'itérations de l'algorithme) et s'il veut que l'on affiche dans la console l'état de l'automate à chaque étape (plus intéressant plutôt que d'avoir l'état à la fin puisqu'entre chaque itération, il y a de l'aléatoire si une personne est contaminée ou non).

### explication des fichiers:
- affichage.go: gère les différentes fonctions pour afficher dans la console ou prodrie l'image de fin

- cellule.go: gère les modifications au niveau d'une cellule ou d'un groupe de cellules (batch)
- client.go: gère la connexion du client
- server.go: pareil pour le serveur
- épidémie.go: teste l'affichage (à supprimer)
- test.go: à renommer, contient les fonctions pour tester les performances des différentes itérations du code et de les comparer à un execution séquentielle de l'algorithme (surtout utile pour le rendu)


### Pourquoi l'avoir implémenté avec des goroutines ?
L'algorithme que l'on a implémenté est en O(n²) au départ et nécessite d'appliquer des modifications ou non à chaque case, les goroutines permetent de diviser le travail et de finir plus vite 
Voici ci-dessous un apperçu des performances pour 100 itérations. A noter que dans le temps de traitement parallèle, il y a aussi le temps d'envoi ds données qu'il n'y a pas pour le traitement séquentiel.

| Taille du tableau | Temps de traitement parallèle | Temps de traitement séquentiel |
|-------------------|-------------------------------|--------------------------------|
| 10x10             | 30 ms                         | 81 ms                         |
| 50x50             | 585 ms                      | 3,23 s                         |
| 60x60             | 2.2s                          | 4.36 s                         |
| 70x70             | 880 ms                        | 5,94 s                          |
| 85x85             | 1,1 s                         | 5,6 s                          |
| 100x100           | 4.4 s                         | 11.7 s                          |
| 125x125           | 9.3 s                        | 18.4 s                         |
| 150x150           | 7.07 s                        | 21.7 s                         |
| 200x200           | 25.5 s                        | 46.8 s                         |

On peut noter qu'il y a des temps qui semblent incohérents pour l'execution en parallèle, cela peut venir de plusieurs facteurs mais comme cela tourne sur plusieurs coeurs, cela peut varier beaucoup en fonction de la charge de traitement que gère déjà mon ordinateur mais l'execution parallèle permet d'executer le code 2 à 3 fois plus vite.

### Comment lancer le projet ?

Pour lancer le projet, rien de plus simple ! il suffit de télécharger l'archive du projet, d'ouvrir une console dans le répertoire, d'écrire go build puis de suivre les indications à l'écran