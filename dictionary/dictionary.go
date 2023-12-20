package dictionary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

// Dictionary représente un dictionnaire
type Dictionary struct {
	words  map[string]string
	mu     sync.RWMutex
	addCh  chan addOperation
	delCh  chan string
	readCh chan readOperation
}

// NewDictionary crée une nouvelle instance de dictionnaire
func NewDictionary() *Dictionary {
	d := &Dictionary{
		words:  make(map[string]string),
		addCh:  make(chan addOperation),
		delCh:  make(chan string),
		readCh: make(chan readOperation),
	}
	go d.processOperations()
	return d
}

// AddWord ajoute un mot et sa définition au dictionnaire
func (d *Dictionary) AddWord(word, definition string) {
	d.addCh <- addOperation{word, definition}
}

// GetDefinition retourne la définition d'un mot spécifique
func (d *Dictionary) GetDefinition(word string) (string, bool) {
	respCh := make(chan readResponse)
	d.readCh <- readOperation{word, respCh}
	resp := <-respCh
	return resp.definition, resp.exists
}

// RemoveWord supprime un mot du dictionnaire
func (d *Dictionary) RemoveWord(word string) {
	d.delCh <- word
}

// PrintDictionary affiche tous les mots et leurs définitions dans le dictionnaire
func (d *Dictionary) PrintDictionary() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	fmt.Println("Dictionnaire:")
	for word, definition := range d.words {
		fmt.Printf("%s: %s\n", word, definition)
	}
}

// SaveToFile enregistre les données du dictionnaire dans un fichier JSON
func (d *Dictionary) SaveToFile(filename string) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// LoadFromFile charge les données du dictionnaire depuis un fichier JSON
func (d *Dictionary) LoadFromFile(filename string) error {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	err = json.Unmarshal(fileData, d)
	if err != nil {
		return err
	}
	return nil
}

// processOperations gère les opérations asynchrones
func (d *Dictionary) processOperations() {
	for {
		select {
		case addOp := <-d.addCh:
			d.mu.Lock()
			d.words[addOp.word] = addOp.definition
			d.mu.Unlock()

		case delWord := <-d.delCh:
			d.mu.Lock()
			delete(d.words, delWord)
			d.mu.Unlock()

		case readOp := <-d.readCh:
			d.mu.RLock()
			definition, exists := d.words[readOp.word]
			d.mu.RUnlock()
			readOp.respCh <- readResponse{definition, exists}
		}
	}
}

// addOperation représente une opération d'ajout
type addOperation struct {
	word       string
	definition string
}

// readOperation représente une opération de lecture
type readOperation struct {
	word   string
	respCh chan readResponse
}

// readResponse représente la réponse à une opération de lecture
type readResponse struct {
	definition string
	exists     bool
}
