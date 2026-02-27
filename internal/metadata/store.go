package metadata

import (
	"sync"
)

type PhotoMetadata struct {
	SourcePath   string
	Category     string
	Tags         []string
	Colors       []string
	Description  string
	DateTaken    string
	CameraMake   string
	CameraModel  string
	Lens         string
	ISO          int
	Aperture     string
	ShutterSpeed string
	FocalLength  string
	Width        int
	Height       int
}

type Store struct {
	mu       sync.RWMutex
	metadata map[string]*PhotoMetadata
}

func NewStore() *Store {
	return &Store{
		metadata: make(map[string]*PhotoMetadata),
	}
}

func (s *Store) Get(relPath string) (*PhotoMetadata, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.metadata[relPath]
	return m, ok
}

func (s *Store) Set(relPath string, m *PhotoMetadata) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metadata[relPath] = m
}

func (s *Store) GetAll() map[string]*PhotoMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]*PhotoMetadata, len(s.metadata))
	for k, v := range s.metadata {
		result[k] = v
	}
	return result
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.metadata)
}
