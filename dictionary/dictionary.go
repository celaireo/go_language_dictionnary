package dictionary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

type addOperation struct {
	word       string
	definition string
	respCh     chan<- readResponse
}

type delOperation struct {
	word   string
	respCh chan<- readResponse
}

type readOperation struct {
	word   string
	respCh chan readResponse
}

type readResponse struct {
	definition string
	exists     bool
}

type dictionaryJSON struct {
	Words map[string]string `json:"words"`
}

type Dictionary struct {
	words  map[string]string
	mu     sync.RWMutex
	addCh  chan addOperation
	delCh  chan delOperation
	readCh chan readOperation
	saveCh chan struct{}
}

func NewDictionary() *Dictionary {
	d := &Dictionary{
		words:  make(map[string]string),
		addCh:  make(chan addOperation),
		delCh:  make(chan delOperation),
		readCh: make(chan readOperation),
		saveCh: make(chan struct{}),
	}
	go d.processOperations()
	return d
}

func (d *Dictionary) AddWord(word, definition string) {
	respCh := make(chan readResponse)
	d.addCh <- addOperation{word, definition, respCh}
	<-respCh
}

func (d *Dictionary) RemoveWord(word string) {
	respCh := make(chan readResponse)
	d.delCh <- delOperation{word, respCh}
	<-respCh
}

func (d *Dictionary) GetDefinition(word string) (string, bool) {
	respCh := make(chan readResponse)
	d.readCh <- readOperation{word, respCh}
	resp := <-respCh
	return resp.definition, resp.exists
}

func (d *Dictionary) PrintDictionary() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	fmt.Println("Dictionnaire:")
	for word, definition := range d.words {
		fmt.Printf("%s: %s\n", word, definition)
	}
}

func (d *Dictionary) SaveToFile(filename string) error {
	log.Println("Saving dictionary to file...")
	close(d.saveCh)
	<-d.saveCh
	d.mu.RLock()
	defer d.mu.RUnlock()
	data, err := json.MarshalIndent(dictionaryJSON{Words: d.words}, "", "  ")
	if err != nil {
		log.Println("Error marshalling dictionary:", err)
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
	} else {
		log.Println("Dictionary saved successfully.")
	}
	return err
}

func (d *Dictionary) LoadFromFile(filename string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &dictionaryJSON{Words: d.words})
}

func (d *Dictionary) processOperations() {
	for {
		select {
		case addOp := <-d.addCh:
			d.mu.Lock()
			d.words[addOp.word] = addOp.definition
			d.mu.Unlock()
			close(addOp.respCh)

		case delOp := <-d.delCh:
			d.mu.Lock()
			delete(d.words, delOp.word)
			d.mu.Unlock()
			close(delOp.respCh)

		case readOp := <-d.readCh:
			d.mu.RLock()
			definition, exists := d.words[readOp.word]
			d.mu.RUnlock()
			readOp.respCh <- readResponse{definition, exists}
		}
	}
}
