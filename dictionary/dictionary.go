package dictionary

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"strings"
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


type addRequest struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
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

func (d *Dictionary) ListWords() map[string]string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make(map[string]string)
	for word, definition := range d.words {
		result[word] = definition
	}
	return result
}


// Sauvegarde dans un fichier JSON
func (d *Dictionary) SaveToFile(filename string) error {
	// Envoie une structure vide via le canal saveCh pour notifier la goroutine principale.
	d.saveCh <- struct{}{}

	// Attend la confirmation de la goroutine principale.
	<-d.saveCh

	// code pour sauvegarder dans un fichier JSON
	jsonData, err := json.Marshal(dictionaryJSON{Words: d.words})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}


func (d *Dictionary) processOperations() {
	for {
		select {
		case addOp := <-d.addCh:
			// ... (code pour traiter l'opération d'ajout)
			d.mu.Lock()
			d.words[addOp.word] = addOp.definition
			d.mu.Unlock()
			addOp.respCh <- readResponse{} // Envoie une réponse vide pour notifier la fin de l'opération.
		
		case delOp := <-d.delCh:
			// ... (code pour traiter l'opération de suppression)
			d.mu.Lock()
			delete(d.words, delOp.word)
			d.mu.Unlock()
			delOp.respCh <- readResponse{} // Envoie une réponse vide pour notifier la fin de l'opération.
		
		case readOp := <-d.readCh:
			// ... (code pour traiter l'opération de lecture)
			d.mu.RLock()
			definition, exists := d.words[readOp.word]
			d.mu.RUnlock()
			readOp.respCh <- readResponse{definition, exists} // Envoie la réponse à l'opération de lecture.
		
		case <-d.saveCh:
			// ... (code pour traiter l'opération de sauvegarde)
			d.saveCh <- struct{}{} // Confirme la fin de l'opération de sauvegarde.
		}
	}
}


// Fonction pour gérer la route HTTP "/list"
func (d *Dictionary) handleList(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	response := make(map[string]string)
	for word, definition := range d.words {
		response[word] = definition
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Fonction pour gérer la route HTTP "/get"
func (d *Dictionary) handleGet(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Path[len("/get/"):]
	if word == "" {
		http.Error(w, "Word parameter is required", http.StatusBadRequest)
		return
	}

	definition, exists := d.GetDefinition(word)
	if !exists {
		http.Error(w, "Word not found in the dictionary", http.StatusNotFound)
		return
	}

	response := map[string]string{"word": word, "definition": definition}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}


// Fonction pour gérer la route HTTP "/add"
func (d *Dictionary) handleAdd(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode de la requête est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lire le corps de la requête
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Désérialiser les données JSON
	var addRequest addRequest
	err = json.Unmarshal(body, &addRequest)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	// Vérifier que les champs nécessaires sont présents
	if addRequest.Word == "" || addRequest.Definition == "" {
		http.Error(w, "Word and definition are required fields", http.StatusBadRequest)
		return
	}

	// Ajouter le mot au dictionnaire
	d.AddWord(addRequest.Word, addRequest.Definition)

	// Répondre avec un statut de succès
	w.WriteHeader(http.StatusCreated)
}

// Fonction pour gérer la route HTTP "/remove"
func (d *Dictionary) handleRemove(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode de la requête est DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire le mot à supprimer de l'URL
	word := r.URL.Path[len("/remove/"):]
	if word == "" {
		http.Error(w, "Word parameter is required", http.StatusBadRequest)
		return
	}

	// Supprimer le mot du dictionnaire
	d.RemoveWord(word)

	// Répondre avec un statut de succès
	w.WriteHeader(http.StatusOK)
}

// Gestion des routes HTTP
func (d *Dictionary) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/list" {
			d.handleList(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/get/") {
			d.handleGet(w, r)
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		d.handleAdd(w, r)
	case http.MethodDelete:
		d.handleRemove(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	
}

	


