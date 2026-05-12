package controller

// controller が「serviceに何をしてほしいか」を定義
type PasteService interface {
	CreatePaste(content string) (string, error)
	GetPaste(id string) (string, error)
}

type PasteController struct {
	service PasteService // ← interfaceで受け取る
}
