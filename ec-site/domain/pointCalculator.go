package domain

type PointCalculator interface {
	CalculatePoints() float32
}

type RegularPointCalculator struct {
	Order *Order
}

type PremiumPointCalculator struct {
	Order *Order
}

// ポイント計算
func (r *RegularPointCalculator) CalculatePoints() float32 {
	return float32(r.Order.TotalPrice) * 0.1
}

func (p *PremiumPointCalculator) CalculatePoints() float32 {
	return float32(p.Order.TotalPrice) * 0.2
}
