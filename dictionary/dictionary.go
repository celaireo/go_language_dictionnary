package dictionary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"os"
)

// Dictionary représente un dictionnaire
type Dictionary struct {
	Words map[string]string `json:"words"`
}

// NewDictionary crée une nouvelle instance de dictionnaire
func NewDictionary() *Dictionary {
	return &Dictionary{
		Words: make(map[string]string),
	}
}

// AddWord ajoute un mot et sa définition au dictionnaire
func (d *Dictionary) AddWord(word, definition string) {
	d.Words[word] = definition
}

// GetDefinition retourne la définition d'un mot spécifique
func (d *Dictionary) GetDefinition(word string) (string, bool) {
	definition, exists := d.Words[word]
	return definition, exists
}

// RemoveWord supprime un mot du dictionnaire
func (d *Dictionary) RemoveWord(word string) {
	delete(d.Words, word)
}

// PrintDictionary affiche tous les mots et leurs définitions dans le dictionnaire
func (d *Dictionary) PrintDictionary() {
	fmt.Println("Dictionnaire:")
	for word, definition := range d.Words {
		fmt.Printf("%s: %s\n", word, definition)
	}
}

// SaveToFile enregistre les données du dictionnaire dans un fichier JSON
func (d *Dictionary) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// LoadFromFile charge les données du dictionnaire depuis un fichier JSON
func (d *Dictionary) LoadFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, d)
}
