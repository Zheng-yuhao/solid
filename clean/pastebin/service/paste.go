package service

import (
	"clean_pastebin/domain"
	"errors"
)

type PasteService struct {
	repo domain.PasteRepository
}

func (s *PasteService) CreatePaste(content string) (string, error) {
	paste := domain.Paste{
		ID:      "123", // Assume this function generates a unique ID
		Content: content,
	}
	if !paste.IsContentValid() {
		return "", errors.New("error")
	}
	err := s.repo.Save(paste)
	if err != nil {
		return "", err
	}
	return paste.ID, nil
}

func (s *PasteService) GetPaste(id string) (string, error) {
	paste, err := s.repo.Get(id)
	if err != nil {
		return "", err
	}
	return paste.Content, nil
}
