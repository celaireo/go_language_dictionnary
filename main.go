package main

import (
	"fmt"
)

func main() {
	// Créer une carte pour stocker les villes et les pays
	dico_villes := make(map[string]string)

	// Ajouter des éléments à la carte
	dico_villes["paris"] = "France"
	dico_villes["lisbonne"] = "Portugal"
	dico_villes["londres"] = "Angleterre"

	// Afficher la carte
	fmt.Println("Carte après ajout:", dico_villes)

	// Utiliser la méthode Get pour obtenir la définition d'un mot spécifique
	if definition, exists := dico_villes["paris"]; exists {
		fmt.Println("Définition de Paris:", definition)
	} else {
		fmt.Println("Paris non trouvé dans le dictionnaire.")
	}

	// Supprimer un élément de la carte
	delete(dico_villes, "lisbonne")

	// Afficher la carte après suppression
	fmt.Println("Carte après suppression:", dico_villes)

	// Lister les éléments de la carte
	fmt.Println("Liste des éléments de la carte:")
	for ville, pays := range dico_villes {
		fmt.Printf("%s: %s\n", ville, pays)
	}
}
