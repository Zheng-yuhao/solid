package domain

// [SRP準拠] Order構造体は注文データの保持のみを担当しており、
// メール送信・決済・ポイント計算などのロジックを持たない。
// 単一の責任(データ表現)に絞られている。
type Order struct {
	Item       []string
	TotalPrice int
	UserID     string
}

func NewOrder(item []string, totalPrice int, userID string) *Order {
	return &Order{
		Item:       item,
		TotalPrice: totalPrice,
		UserID:     userID,
	}
}
