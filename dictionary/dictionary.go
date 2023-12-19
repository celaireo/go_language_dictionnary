
package dictionary

import "fmt"

// Dictionary représente un dictionnaire
type Dictionary struct {
	words map[string]string
}

// NewDictionary crée une nouvelle instance de dictionnaire
func NewDictionary() *Dictionary {
	return &Dictionary{
		words: make(map[string]string),
	}
}

// AddWord ajoute un mot et sa définition au dictionnaire
func (d *Dictionary) AddWord(word, definition string) {
	d.words[word] = definition
}

// GetDefinition retourne la définition d'un mot spécifique
func (d *Dictionary) GetDefinition(word string) (string, bool) {
	definition, exists := d.words[word]
	return definition, exists
}

// RemoveWord supprime un mot du dictionnaire
func (d *Dictionary) RemoveWord(word string) {
	delete(d.words, word)
}

// PrintDictionary affiche tous les mots et leurs définitions dans le dictionnaire
func (d *Dictionary) PrintDictionary() {
	fmt.Println("Dictionnaire:")
	for word, definition := range d.words {
		fmt.Printf("%s: %s\n", word, definition)
	}
}




// package main

// import (
// 	"fmt"
// )

// func main() {
// 	// Créer une carte pour stocker les villes et les pays
// 	dico_villes := make(map[string]string)

// 	// Ajouter des éléments à la carte
// 	addToDictionary(dico_villes, "paris", "France")
// 	addToDictionary(dico_villes, "lisbonne", "Portugal")
// 	addToDictionary(dico_villes, "londres", "Angleterre")

// 	// Afficher la carte
// 	fmt.Println("Carte après ajout:", dico_villes)

// 	// Utiliser la méthode Get pour obtenir la définition d'un mot spécifique
// 	definition, exists := getFromDictionary(dico_villes, "paris")
// 	if exists {
// 		fmt.Println("Définition de Paris:", definition)
// 	} else {
// 		fmt.Println("Paris non trouvé dans le dictionnaire.")
// 	}

// 	// Supprimer un élément de la carte
// 	removeFromDictionary(dico_villes, "lisbonne")

// 	// Afficher la carte après suppression
// 	fmt.Println("Carte après suppression:", dico_villes)

// 	// Lister les éléments de la carte
// 	listDictionary(dico_villes)
// }

// // Fonction pour ajouter un élément à la carte
// func addToDictionary(dico map[string]string, ville string, pays string) {
// 	dico[ville] = pays
// }

// // Fonction pour obtenir la définition d'un mot de la carte
// func getFromDictionary(dico map[string]string, ville string) (string, bool) {
// 	definition, exists := dico[ville]
// 	return definition, exists
// }

// // Fonction pour supprimer un élément de la carte
// func removeFromDictionary(dico map[string]string, ville string) {
// 	delete(dico, ville)
// }

// // Fonction pour lister les éléments de la carte
// func listDictionary(dico map[string]string) {
// 	fmt.Println("Liste des éléments de la carte:")
// 	for ville, pays := range dico {
// 		fmt.Printf("%s: %s\n", ville, pays)
// 	}
// }
