package main

import (
	"fmt"
	"dictionnaire/dictionary"
	"sync"
)

func main() {
	// Créer une nouvelle instance de votre dictionnaire
	dico := dictionary.NewDictionary()

	// Utiliser un WaitGroup pour attendre la fin des opérations asynchrones
	var wg sync.WaitGroup

	// Lancer une goroutine pour l'ajout d'éléments
	wg.Add(1)
	go func() {
		defer wg.Done()
		dico.AddWord("paris", "France")
		dico.AddWord("lisbonne", "Portugal")
		dico.AddWord("londres", "Angleterre")
		dico.AddWord("abidjan", "Côte d'Ivoire")
	}()

	// Lancer une goroutine pour la suppression d'un élément
	wg.Add(1)
	go func() {
		defer wg.Done()
		dico.RemoveWord("lisbonne")
	}()

	// Lancer une goroutine pour l'affichage du dictionnaire
	wg.Add(1)
	go func() {
		defer wg.Done()
		dico.PrintDictionary()
	}()

	// Attendre la fin de toutes les opérations asynchrones
	wg.Wait()

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
