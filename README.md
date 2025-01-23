repo pour GO en tc

Le projet consiste en un automate cellulaire pour simuler la propagation d'une épidémie dans une population.

On laisse l'utilisateur décider de la durée de la simulation (le nombre d'itérations de l'algorithme) et s'il veut que l'on affiche dans la console l'état de l'automate à chaque étape (plus intéressant plutôt que d'avoir l'état à la fin puisqu'entre chaque itération, il y a de l'aléatoire si une personne est contaminée ou non).

explication des fichiers:
affichage.go: gère les différentes fonctions pour afficher dans la console ou prodrie l'image de fin

cellule.go: gère les modifications au niveau d'une cellule ou d'un groupe de cellules (batch)
client.go: gère la connexion du client
server.go: pareil pour le serveur
épidémie.go: teste l'affichage (à supprimer)
brouillon_workerpool.go: quelques idées pour implémenter un workerpool (à supprimer plus tard)
test.go: à renommer, contient les fonctions pour tester les performances des différentes itérations du code et de les comparer à un execution séquentielle de l'algorithme (surtout utile pour le rendu)


Pourquoi l'avoir implémenté avec des goroutines ?
L'algorithme que l'on a implémenté est en O(n²) au départ et nécessite d'appliquer des modifications ou non à chaque case, les goroutines permetent de diviser le travail et de finir plus vite 