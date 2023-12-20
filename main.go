package main

import (
	"dictionnaire/dictionary"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	dico := dictionary.NewDictionary()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		dico.AddWord("paris", "France")
		dico.AddWord("lisbonne", "Portugal")
		dico.AddWord("londres", "Angleterre")
		dico.AddWord("abidjan", "Côte d'Ivoire")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dico.RemoveWord("lisbonne")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dico.PrintDictionary()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		err := dico.SaveToFile("dictionary.json")
		if err != nil {
			log.Println("Error saving dictionary to file:", err)
		}
	}()

	wg.Wait()

	time.Sleep(1 * time.Second)

	definition, exists := dico.GetDefinition("paris")
	if exists {
		fmt.Println("Définition de Paris:", definition)
	} else {
		fmt.Println("Paris non trouvé dans le dictionnaire.")
	}

	err := dico.LoadFromFile("dictionary.json")
	if err != nil {
		log.Println("Erreur lors du chargement du dictionnaire :", err)
		return
	}

	fmt.Println("Dictionnaire après chargement depuis le fichier:")
	dico.PrintDictionary()
}
