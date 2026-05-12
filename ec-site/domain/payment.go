package domain

import (
	"fmt"
)

// [ISP準拠] PaymentMethodインターフェースはProcess一つのみを定義しており、
// 不必要なメソッドを実装側に強制していない。
//
// [OCP準拠] 新しい支払い方法を追加する場合、このインターフェースを実装した
// 新しい構造体を追加するだけでよく、既存のコードを変更する必要がない。
type PaymentMethod interface {
	// [Bug修正] process → Process (メソッド名をexportedに変更)
	// 原因: 小文字始まりはunexportedのため、domainパッケージ外から呼び出せずコンパイルエラーになる。
	Process(*Order) error
}

// [LSP準拠] CreditCardPaymentはPaymentMethodを実装しており、
// PaymentMethodが期待される箇所でそのまま代替できる。
type CreditCardPayment struct {
	Order *Order
}

// [LSP準拠] BankTransferPaymentもPaymentMethodを実装しており、
// CreditCardPaymentと同様にPaymentMethodとして代替可能。
type BankTransferPayment struct {
	Order *Order
}

func (p *CreditCardPayment) Process(order *Order) error {
	fmt.Println("クレジットカード処理")
	return nil
}

func (p *BankTransferPayment) Process(order *Order) error {
	fmt.Println("銀行振込処理")
	return nil
}

// [OCP違反] NewPaymentMethodのswitch文は、新しい支払い方法を追加するたびに
// この関数を修正する必要があり、OCPに違反している。
// Mapや登録パターン(registry pattern)を使うことで解消できる。→簡便の多分ここは一旦無視でええ
func NewPaymentMethod(method string, order *Order) PaymentMethod {
	switch method {
	case "credit_card":
		return &CreditCardPayment{Order: order}
	case "bank_transfer":
		return &BankTransferPayment{Order: order}
	default:
		return nil
	}
}
