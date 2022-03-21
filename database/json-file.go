package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/pkg/errors"
)

// SafeJsonFile represents a Json file on the filesystem. It is protected by a mutex, it can be called concurrently.
// Can be used as poor man's database.
type SafeJsonFile struct {
	mu   sync.Mutex
	path string
}

// NewSafeJsonFile creates a new SafeJsonFile object with the supplied path. If the file not already exists,
// it will be created. NewSafeJsonFile may return an error if the file cannot be opened or created.
func NewSafeJsonFile(path string) (*SafeJsonFile, error) {
	if jsonFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return nil, errors.Wrap(err, "Failed to ensure Json DB file")
	} else {
		if err := jsonFile.Close(); err != nil {
			return nil, errors.Wrap(err, "Failed to ensure Json DB file")
		}
	}

	return &SafeJsonFile{
		path: path,
	}, nil
}

// Load load loads the content of the file into the supplied data object.
func (sjf *SafeJsonFile) Load(data interface{}) error {
	sjf.mu.Lock()
	defer sjf.mu.Unlock()

	jsonBytes, err := ioutil.ReadFile(sjf.path)
	if err != nil {
		return errors.Wrap(err, "Failed to read DB Json")
	}

	return json.Unmarshal([]byte(jsonBytes), data)
}

// Save saves the supplied data object into the file.
func (sjf *SafeJsonFile) Save(data interface{}) error {
	sjf.mu.Lock()
	defer sjf.mu.Unlock()

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Failed to encode data as Json")
	}

	return ioutil.WriteFile(sjf.path, jsonBytes, 0644)
}

// Ping checks if the Json file is still available.
func (sjf *SafeJsonFile) Ping() error {
	_, err := os.Stat(sjf.path)
	if err != nil {
		return errors.Wrap(err, "Failed to check Json DB")
	}
	return nil
}
