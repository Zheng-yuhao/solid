# Clean Architecture 実装ガイド

## 唯一のルール

> **依存の向きは常に内向き。内側は外側を知らない。**

層の数・名前はプロジェクトごとに変えてよい。このルールだけは変えない。

---

## 層の構成テンプレート

```
┌─────────────────────────────────┐
│  最外層（HTTP / CLI / gRPC）     │  フレームワーク・入出力
├─────────────────────────────────┤
│  adapter 層（controller 等）     │  外の世界 ↔ 内側の変換
├─────────────────────────────────┤
│  service 層（use case）          │  ビジネス手順
├─────────────────────────────────┤
│  domain 層                      │  業務ルール・不変の核
└─────────────────────────────────┘
         ▲ 実装する
┌─────────────────────────────────┐
│  repository 層                  │  DB・外部 API の具体的な実装
└─────────────────────────────────┘
```

層を増やすときは、追加した層が **内側の interface に依存しているか** だけ確認する。

---

## 各層の責任

| 層 | 書くもの | 書かないもの |
|---|---|---|
| domain | struct・業務ルールメソッド・Repository interface | SQL・HTTP・フレームワーク |
| service | ビジネス手順・バリデーション呼び出し | SQL・HTTPステータスコード |
| controller | HTTP パース・レスポンス変換・Service interface | SQL・業務ルール |
| repository | SQL・外部 API 呼び出し | 業務ルール・HTTPの知識 |

---

## interface をどこに置くか

```
「domain が必要としているか？」   → domain に置く   （Repository interface）
「使う側が要求しているか？」      → 使う側に置く    （Service interface 等）
```

| interface | 置く場所 | 実装する場所 |
|---|---|---|
| `XxxRepository` | domain | repository |
| `XxxService` | controller | service |

---

## 依存の向き（import で確認する）

```
controller  →  domain
service     →  domain
repository  →  domain

domain  →  誰も import しない ← これが守れていれば正しい
```

`go build` で循環 import エラーが出たら、内側が外側を import しているサイン。

---

## 新機能追加の手順

```
① domain
     struct を定義する
     業務ルールメソッドを書く（DB・HTTP に無関係な判断）
     XxxRepository interface を定義する

② repository
     domain.XxxRepository interface を満たす struct を実装する
     SQL・外部 API の呼び出しはここにだけ書く

③ service
     domain.XxxRepository を field に持つ struct を定義する
     バリデーション → domain 操作 → repo 呼び出しの順で書く

④ controller
     XxxService interface を定義する（自分が使うメソッドだけ）
     HTTP パース → service 呼び出し → レスポンス返却を書く

⑤ テスト
     各層のテストファイルに Mock を定義して単体テストを書く
```

---

## テストと Mock

### 考え方

```
テスト対象の層が持つ interface → Mock に差し替える → DB・HTTP なしでテストできる
```

interface に依存しているから Mock に差し替えられる。  
具体的な struct に依存していたら差し替えられない。

### Mock の配置

| テスト対象 | Mock を定義する場所 | 差し替えるもの |
|---|---|---|
| domain | 不要 | 外部依存なし |
| service | `service/xxx_test.go` | `domain.XxxRepository` |
| controller | `controller/xxx_test.go` | `XxxService` interface |

### Mock の書き方テンプレート

```go
// service のテスト: domain.XxxRepository を差し替える
type mockXxxRepository struct {
    store      map[string]domain.Xxx
    saveCalled bool
}

func (m *mockXxxRepository) Save(xxx domain.Xxx) error {
    m.store[xxx.ID] = xxx
    m.saveCalled = true
    return nil
}

func (m *mockXxxRepository) Get(id string) (domain.Xxx, error) {
    xxx, ok := m.store[id]
    if !ok {
        return domain.Xxx{}, errors.New("not found")
    }
    return xxx, nil
}

func TestXxxService_Create_正常系(t *testing.T) {
    // Arrange
    repo := &mockXxxRepository{store: make(map[string]domain.Xxx)}
    svc := &XxxService{repo: repo}

    // Act
    id, err := svc.Create("input")

    // Assert
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if !repo.saveCalled {
        t.Fatal("expected Save to be called")
    }
    _ = id
}
```

```go
// controller のテスト: XxxService interface を差し替える
type mockXxxService struct {
    createResult string
    createErr    error
}

func (m *mockXxxService) Create(input string) (string, error) {
    return m.createResult, m.createErr
}

func TestXxxController_Create_正常系(t *testing.T) {
    // Arrange
    svc := &mockXxxService{createResult: "new-id"}
    ctrl := XxxController{service: svc}

    req := httptest.NewRequest(http.MethodPost, "/xxx", strings.NewReader(`{"input":"hello"}`))
    w := httptest.NewRecorder()

    // Act
    ctrl.Create(w, req)

    // Assert
    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }
}
```

---

## 実装チェックリスト

### domain

- [ ] struct のフィールドはすべて業務上意味のある項目か
- [ ] メソッドは DB・HTTP と無関係な判断・操作だけか
- [ ] `XxxRepository` interface が定義されているか
- [ ] domain が他のパッケージを import していないか

### repository

- [ ] `domain.XxxRepository` interface をすべて実装しているか（シグネチャが一致しているか）
- [ ] `domain.Xxx` 型を引数・戻り値に使っているか（独自型になっていないか）
- [ ] SQL・外部 API の呼び出しがこの層だけに閉じているか
- [ ] 業務ルールの判断を書いていないか

### service

- [ ] `repo` フィールドの型が具体的な struct ではなく interface か
- [ ] バリデーションを repo 呼び出しより前にしているか
- [ ] HTTP の知識（ステータスコードなど）を持ち込んでいないか
- [ ] constructor（`NewXxxService`）が interface 型を返しているか

### controller

- [ ] `XxxService` interface を自パッケージ内に定義しているか
- [ ] `service` フィールドの型が interface か
- [ ] HTTP パース・バリデーション・レスポンス変換だけに集中しているか
- [ ] 業務ルールの判断を書いていないか

### テスト

- [ ] service のテストで DB を使っていないか（Mock を使っているか）
- [ ] controller のテストで実 service を使っていないか（Mock を使っているか）
- [ ] 正常系・異常系（空入力・not found・エラー返却）の両方があるか
- [ ] Mock に `saveCalled` 等のフラグを持たせて「呼ばれたか」を検証しているか

### 全体

- [ ] `go build ./...` が通るか
- [ ] 循環 import が発生していないか
- [ ] domain が誰も import していないか
- [ ] 各層が「自分より内側の interface だけ」に依存しているか
