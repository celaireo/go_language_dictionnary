package main

import (
	"fmt"
	"dictionnaire/dictionary"
)

func main() {
	// Créer une nouvelle instance de votre dictionnaire
	dico := dictionary.NewDictionary()

	// Ajouter des éléments au dictionnaire
	dico.AddWord("paris", "France")
	dico.AddWord("lisbonne", "Portugal")
	dico.AddWord("londres", "Angleterre")
	dico.AddWord("abidjan", "Côte d'Ivoire")

	// Afficher le dictionnaire
	fmt.Println("Dictionnaire après ajout:")
	dico.PrintDictionary()

	// Enregistrer le dictionnaire dans un fichier JSON
	err := dico.SaveToFile("dictionary.json")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement du dictionnaire :", err)
		return
	}

	// Utiliser la méthode Get pour obtenir la définition d'un mot spécifique
	definition, exists := dico.GetDefinition("paris")
	if exists {
		fmt.Println("Définition de Paris:", definition)
	} else {
		fmt.Println("Paris non trouvé dans le dictionnaire.")
	}

	// Supprimer un mot du dictionnaire
	dico.RemoveWord("lisbonne")

	// Afficher le dictionnaire après suppression
	fmt.Println("Dictionnaire après suppression:")
	dico.PrintDictionary()

	// Charger le dictionnaire depuis le fichier JSON
	err = dico.LoadFromFile("dictionary.json")
	if err != nil {
		fmt.Println("Erreur lors du chargement du dictionnaire :", err)
		return
	}

	// Afficher le dictionnaire après chargement depuis le fichier
	fmt.Println("Dictionnaire après chargement depuis le fichier:")
	dico.PrintDictionary()
}
