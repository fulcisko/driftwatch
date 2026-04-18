package drift

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Annotation struct {
	Service   string    `json:"service"`
	Key       string    `json:"key"`
	Note      string    `json:"note"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type annotationStore struct {
	Annotations []Annotation `json:"annotations"`
}

func AddAnnotation(path, service, key, note, author string) error {
	if service == "" || key == "" || note == "" {
		return errors.New("service, key, and note are required")
	}
	store, _ := loadAnnotations(path)
	store.Annotations = append(store.Annotations, Annotation{
		Service:   service,
		Key:       key,
		Note:      note,
		Author:    author,
		CreatedAt: time.Now().UTC(),
	})
	return saveAnnotations(path, store)
}

func LoadAnnotations(path string) ([]Annotation, error) {
	store, err := loadAnnotations(path)
	return store.Annotations, err
}

func FilterAnnotations(annotations []Annotation, service, key string) []Annotation {
	var out []Annotation
	for _, a := range annotations {
		if (service == "" || a.Service == service) && (key == "" || a.Key == key) {
			out = append(out, a)
		}
	}
	return out
}

func loadAnnotations(path string) (annotationStore, error) {
	var store annotationStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveAnnotations(path string, store annotationStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
