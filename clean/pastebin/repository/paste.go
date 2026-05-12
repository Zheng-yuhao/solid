package repository

import (
	"clean_pastebin/domain"
)

type Paste struct {
	memoryDB map[string]domain.Paste
}

func (p *Paste) Save(paste domain.Paste) error {
	// Save the paste to the database
	return nil
}

func (p *Paste) Get(id string) (domain.Paste, error) {
	// Retrieve the paste from the database using the id
	return domain.Paste{}, nil
}
