# SOLID / Clean Architecture Projects

SOLID原則と Clean Architecture を Go で実践的に学ぶためのサンプルプロジェクト集。

参考: [SOLID原則の実践的ケーススタディ](https://shuji-bonji.github.io/Notes-on-SOLID-Principle/practical-case-studies.html)

---

## ディレクトリ構成

```
solid/
├── ec-site/        # SOLID原則の基礎演習（単一ドメイン）
├── pastebin/       # SOLID原則の応用演習（レイヤードアーキテクチャ）
└── clean/
    └── pastebin/   # Clean Architecture の実装練習
```

---

## プロジェクト概要

| | ec-site | pastebin | clean/pastebin |
|---|---|---|---|
| 目的 | SOLID原則の基礎を一通り体験する | より実践的なアーキテクチャでSOLIDを適用する | Clean Architecture の層構造と依存の向きを体得する |
| 複雑さ | シンプル（単一ドメイン） | 中規模（レイヤードアーキテクチャ） | 中規模（4層構成） |
| 主な学習テーマ | SRP / DIP / LSP / ISP の基本 | OCP（レジストリパターン）/ DIP（Goのinterface慣習） | 層の責任分離・interface の置き場所・Mock テスト |
| バグ修正あり | あり（意図的なバグを修正する演習） | なし（LSP違反を意図的に残している） | なし |

---

## ec-site

ECサイトの注文処理を題材にしたシンプルなSOLID演習。バグ修正とSOLIDレビューをセットで学ぶ。

### ファイル構成

```
ec-site/
├── app/
│   └── main.go             # エントリポイント・依存注入
└── domain/
    ├── order.go            # 注文データ
    ├── payment.go          # 決済インターフェース・実装
    ├── emailService.go     # メール送信インターフェース・実装
    └── pointCalculator.go  # ポイント計算インターフェース・実装
```

### SOLID対応状況

#### ✅ SRP（単一責任の原則）
- `Order` はデータ保持のみ。メール・決済・ポイント計算のロジックを一切持たない
- `EmailService`・`CreditCardPayment`・`RegularPointCalculator` はそれぞれ1つの責任に絞られている

#### ✅ DIP（依存性逆転の原則）
- `orderService` は `EmailSender`・`PaymentMethod`・`PointCalculator` の各インターフェースに依存し、具体型には依存しない
- `main.go` で具体型を生成し、インターフェースとして注入している

#### ✅ LSP（リスコフの置換原則）
- `CreditCardPayment` と `BankTransferPayment` は共に `PaymentMethod` を実装しており、相互に代替可能
- `RegularPointCalculator` と `PremiumPointCalculator` も `PointCalculator` を通じて代替可能

#### ✅ ISP（インターフェース分離の原則）
- `PaymentMethod` は `Process` のみ、`EmailSender` は `SendConfirmationEmail` のみと、インターフェースが最小化されている

#### ⚠️ OCP（開放閉鎖の原則）違反（意図的に残存）
- `NewPaymentMethod` の `switch` 文は支払い方法追加のたびに修正が必要
- 学習の簡便化のため意図的に未修正

### 修正したバグ一覧

| ファイル | バグ内容 | 原因 |
|---|---|---|
| `main.go` | `package app` | Goのエントリポイントは `package main` 必須 |
| `main.go` | `*domain.PaymentMethod`（ポインタ） | interfaceはすでに参照型なのでポインタ不要 |
| `main.go` | `newOrderServiice`（typo） | メソッド名のタイプミス |
| `main.go` | `paymentSVS`（typo） | フィールド名のタイプミス |
| `payment.go` | `process`（小文字） | unexportedのためパッケージ外から呼べない |
| `emailService.go` | `e.order`（小文字） | exportedフィールドは `e.Order` でアクセスする |
| `emailService.go` | `sendConfirmationEmail`（小文字） | unexportedのためパッケージ外から呼べない |

---

## pastebin

Pastebinサービスを題材にした、レイヤードアーキテクチャによるSOLID実践。

### ファイル構成

```
pastebin/
├── app/
│   └── main.go                         # エントリポイント
├── domain/
│   └── user.go                         # ユーザードメインモデル
├── paste/
│   └── service.go                      # ShortLinkGenerator interface・Factoryとregistry
├── user/
│   └── service.go                      # UserService・UserRepository interface
├── internal/
│   └── repository/
│       └── mysql/
│           ├── user.go                 # UserRepository実装
│           ├── paste.go                # (スケルトン)
│           └── comment.go              # (スケルトン)
└── short_link_generator/
    ├── md5.go                          # MD5実装（init()で自動登録）
    └── sha256.go                       # SHA256実装（init()で自動登録）
```

### SOLID対応状況

#### ✅ SRP（単一責任の原則）
- `domain.User` はユーザーデータの保持のみ
- `UserService` はユーザー操作のユースケース層のみ
- `mysql.UserRepository` はDB操作のみ
- 各層が明確に分離されている

#### ✅ OCP（開放閉鎖の原則）
- `paste.Register` によるレジストリパターンで、新しい短縮URLアルゴリズム追加時に既存コードを変更不要
- `md5.go` と `sha256.go` は `init()` で自動登録されるため、`Factory` 側を一切触らずに拡張できる

```go
// 新アルゴリズムを追加する場合：新ファイルを作るだけでよい
func init() {
    paste.Register("base62", func(baseUrl string) paste.ShortLinkGenerator {
        return &base62Generator{baseUrl}
    })
}
```

#### ✅ DIP（依存性逆転の原則）
- `UserRepository` インターフェースは `user` パッケージ（上位モジュール＝使う側）に定義されており、Goの慣習「interfaceは使う側が定義する」に従っている
- `ShortLinkGenerator` インターフェースも `paste` パッケージ（上位）が所有

#### ✅ ISP（インターフェース分離の原則）
- `UserRepository` は `GetByID` と `CreateUser` のみで最小限
- `ShortLinkGenerator` は `GenerateShortLink` のみ

#### ✅ LSP（リスコフの置換原則）
- `md5Generator` と `sha256Generator` はどちらも `ShortLinkGenerator` の事後条件（有効なURLを返す）を満たすべき実装
- `ShortLinkGeneratorFactory` が未登録アルゴリズム時に `nil` を返す点は**意図的なLSP違反として残している**（呼び出し元はnilを期待しないためnilを返すと事後条件違反になる。`error` を返す設計が正しい）

---

---

## clean/pastebin

Pastebin サービスを題材にした Clean Architecture の実装練習。  
詳細は [`clean/pastebin/README.md`](clean/pastebin/README.md) を参照。

### ファイル構成

```
clean/pastebin/
├── domain/         # 業務ルール・struct・Repository interface
├── service/        # ビジネス手順（use case）
├── repository/     # DB の具体的な実装（メモリ）
└── controller/     # HTTP の橋渡し・Service interface
```

### 層の依存関係

```
controller ──▶ domain （Service interface 経由）
service    ──▶ domain （Repository interface 経由）
repository ──▶ domain （Repository interface を実装）

domain は誰も import しない
```

### 学習ポイント

#### ✅ 依存の向き
- 内側（domain）は外側を知らない
- 外側が内側の interface に依存することで差し替え可能になる

#### ✅ interface の置き場所
- `PasteRepository` → domain（domain が「自分に必要なもの」を宣言）
- `PasteService` → controller（controller が「使いたいもの」を宣言）

#### ✅ Mock によるテスト分離
- service テストでは `domain.PasteRepository` を Mock に差し替えて DB 不要でテストできる
- controller テストでは `PasteService` interface を Mock に差し替えて service 不要でテストできる

---

## 学習の使い分け

| 学びたいこと | 参照先 |
|---|---|
| SOLID各原則の基本概念 | ec-site |
| バグとSOLID違反の関係を体感する | ec-site（バグ修正演習） |
| OCPをレジストリパターンで実践する | pastebin / `paste/service.go` + `short_link_generator/` |
| GoのDIP慣習（interfaceは使う側が定義） | pastebin / `user/service.go` |
| LSP違反の具体例（事後条件） | pastebin / `paste/service.go` の `ShortLinkGeneratorFactory` |
| レイヤードアーキテクチャとSOLIDの関係 | pastebin 全体 |
| Clean Architecture の層構造と依存の向き | clean/pastebin 全体 |
| interface の置き場所の判断基準 | clean/pastebin / `domain/paste.go` + `controller/paste.go` |
| Mock を使った単体テストの書き方 | clean/pastebin / `README.md` のテストセクション |
