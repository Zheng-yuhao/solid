package domain

type Paste struct {
	ID      string
	Content string
}

type PasteRepository interface {
	Save(paste Paste) error
	Get(id string) (Paste, error)
}

func (p *Paste) IsContentValid() bool {
	if p.Content == "" {
		return false
	}
	return true
}
