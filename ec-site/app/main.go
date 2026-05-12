/*
 * 本PJはhttps://shuji-bonji.github.io/Notes-on-SOLID-Principle/practical-case-studies.html
 * を元に、SOLID原則を適用してECサイトのコードを実装するものです。
 */

// [Bug修正] package app → package main
// 原因: Goではエントリポイントのfunc main()は必ずpackage mainに存在しなければならない。
// package appのままではビルド時に「main関数が見つからない」エラーになる。
package main

import (
	"ec-site/domain"
)

// [DIP違反] emailSVCフィールドが具体型 *domain.EmailService に依存している。
// DIPでは上位モジュール(orderService)は抽象(interface)に依存すべき。
// EmailServiceのinterfaceが定義されていないため、差し替えやテストが困難になる。
//
// [SRP準拠] orderServiceは注文処理の調整(orchestration)のみを担当しており、
// メール送信・決済・ポイント計算の具体的な実装を持たない点はSRPに沿っている。
type orderService struct {
	emailSVC domain.EmailSender
	// [Bug修正] *domain.PaymentMethod → domain.PaymentMethod
	// 原因: PaymentMethodはinterfaceであり、Goではinterfaceへのポインタ(*interface)は使用しない。
	// interfaceはすでに参照型なのでポインタ不要。
	//
	// [DIP準拠] PaymentMethodインターフェースに依存しており、具体的な実装(CreditCard等)に依存していない。
	paymentSVC domain.PaymentMethod
	// [Bug修正] *domain.PointCalculator → domain.PointCalculator
	// 原因: PaymentMethodと同様、PointCalculatorもinterfaceなのでポインタ不要。
	//
	// [DIP準拠] PointCalculatorインターフェースに依存している。
	pointCalculator domain.PointCalculator
}

func newOrderService(emailSVC domain.EmailSender, paymentSVC domain.PaymentMethod, pointCalculator domain.PointCalculator) *orderService {
	return &orderService{
		emailSVC:        emailSVC,
		paymentSVC:      paymentSVC,
		pointCalculator: pointCalculator,
	}
}

func main() {
	order := domain.NewOrder([]string{"item1", "item2"}, 1000, "user123")
	emailSVC := &domain.EmailService{Order: order}
	paymentSVC := domain.NewPaymentMethod("credit_card", order)
	pointCalculator := &domain.RegularPointCalculator{Order: order}

	// [Bug修正] newOrderServiice → newOrderService (typo修正)
	// [Bug修正] &paymentSVS → paymentSVC (フィールド名typo、かつinterfaceへのポインタ不要)
	orderSVC := newOrderService(emailSVC, paymentSVC, pointCalculator)

	// [Bug修正] orderSVC.paymentSVS → orderSVC.paymentSVC
	// 原因: フィールド名はpaymentSVCだがpaymentSVSと誤記していた。
	orderSVC.paymentSVC.Process(order)
	orderSVC.emailSVC.SendConfirmationEmail()
	points := orderSVC.pointCalculator.CalculatePoints()
	println("獲得ポイント:", points)
}
