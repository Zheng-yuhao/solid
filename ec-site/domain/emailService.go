package domain

import (
	"fmt"
)

// [SRP準拠] EmailServiceはメール送信のみを担当しており、単一の責任を持つ。
//
// [DIP違反] EmailServiceに対応するinterfaceが定義されていない。
// orderServiceが *EmailService という具体型に依存しているため、
// テスト時にモックへの差し替えや別実装への切り替えができない。
// EmailSender等のinterfaceを定義してそれに依存させるべき。

type EmailSender interface {
	SendConfirmationEmail()
}

type EmailService struct {
	Order *Order
}

// [Bug修正] e.order → e.Order
// 原因: Orderフィールドは大文字(exported)で定義されているが、
// e.order(小文字)でアクセスしておりフィールドが存在せずコンパイルエラーになる。
//
// [Bug修正] sendConfirmationEmail → SendConfirmationEmail (メソッド名をexportedに変更)
// 原因: 小文字始まりのメソッドはGoではunexported扱いになり、
// domainパッケージ外(appパッケージ)からは呼び出せずコンパイルエラーになる。
func (e *EmailService) SendConfirmationEmail() {
	fmt.Printf("ユーザーID: %s に確認メールを送信\n", e.Order.UserID)
}
