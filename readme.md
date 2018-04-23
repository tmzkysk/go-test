## はじめに
golangでのテストはとてもシンプルで、rubyのrspecのように新しくDSLを覚える必要もありません。  
テストについて解説しているサイトは沢山あるのですが、自分の中で特にこれは最初に覚えておいた方がいいなと思うことをピックアップしました。

## 基本的なtestの書き方
- 例えばcalc.go のテストならば同じディレクトリ内に`calc_test.go`という名前で作成する。
- テストファイル内では`testing`パッケージをインポートする。
- テストファイル内では、`TestXXX`という名前でテストメソッドを作成する。
- DSLは特に無いので普通にテストコードを書く。
- パラメータと期待値の組み合わせの配列を用意して、ループで検証していく形が推奨されている(Table Driven Test)

```go:calc.go
package calc

func Add(a,b int) int {
	return a + b
}
```

```go:calc_test.go
package calc

import (
	"testing"
)

func TestAdd(t *testing.T) {
	patterns := []struct {
		a        int
		b        int
		expected int
	}{
		{1, 2, 3},
		{10, -2, 8},
		{-10, -2, -12},
	}

	for idx, pattern := range patterns {
		actual := Add(pattern.a, pattern.b)
		if pattern.expected != actual {
			t.Errorf("pattern %d: want %d, actual %d", idx, pattern.expected, actual)
		}
	}
}
```

## testの実行方法
- カレントディレクトリ以下すべてを再帰的にテスト`go test -v ./...`
- 特定のパッケージをテスト`go test -v ./hogehoge`(パッケージディレクトリを相対パスで指定する)
- 特定のメソッドのみテストする`go test -run TestAdd ./...`

※ -v オプションを付けると実行結果に詳細が付きますので、基本的にはつけておいたほうが良いです。

## テストの実行前後に処理を入れるには
TestMainメソッドを定義します。  
`code := m.Run()`を実行するとテストメソッドが走るので、その前後にDBの初期化処理等を入れることが出来ます。

```go:calc_test.go
package calc

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("before test")
	code := m.Run()
	fmt.Println("after test")
	os.Exit(code)
}

func TestAdd(t *testing.T) {
	// 以下省略
}
```

これを実行すると、以下のようになります。  
テストの前後にfmt.Printlnが入っているのがわかります。

```sh
$ go test -v ./...
before test
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
PASS
after test
ok
```

## テストでモックを使うには
### インターフェースを使ったモック
インターフェースを使っているオブジェクトの場合、実際のコードとテストコードでインタフェースに定義するオブジェクトを変えることでテスト時の振る舞いを変えることが出来ます。
ここではsomefuncパッケージのClientオブジェクトのRunメソッド内で呼び出されるcallメソッドの振る舞いを、モックを使って切り替える方法を紹介します。

```go:somefunc.go
package somefunc

type Caller interface {
	call(val int) int
}

type Client struct {
	FuncCaller Caller
}

type ExampleCaller struct{}

func (c *Client) Run(val int) int {
	return c.FuncCaller.call(val)
}

func (f *ExampleCaller) call(val int) int {
	return val
}
```

上記のコードを実行するには以下のように呼び出します。

```go:main.go
c := somefunc.Client{&somefunc.ExampleCaller{}}
c.Run(1)
```

ここで、テスト時にExampleCallerのモックを作って、callメソッドの振る舞いを変えるにはテストコードを以下のようにします。

```go:somefunc_test.go
package somefunc

import (
	"testing"
)

func TestRun(t *testing.T) {

	patterns := []struct {
		val      int
		expected int
	}{
		{2, 2},
		{8, 8},
		{-10, -10},
	}

	for idx, pattern := range patterns {
		// Clientのnewの際に、モックオブジェクトを引数にする
		c := Client{&mockCaller{}}
		actual := c.Run(pattern.val)
		if pattern.expected != actual {
			t.Errorf("pattern %d: want %d, actual %d", idx, pattern.expected, actual)
		}
	}
}

// callメソッドのレシーバをmockCallerとして宣言する。
type mockCaller struct{}

// 通常のコードではcallメソッドは引数の値をそのまま返却するが、
// モックでは、引数 + 10した値を返却するようにする。
func (s *mockCaller) call(val int) int {
	return val + 10
}
```

### 変数の再代入で行う方法
ここではsomeprocessパッケージのRun関数のテストを行っていますが、Run内でcallという関数を呼び出しています。  
このcall関数の挙動をテストの時だけ切り替えるには、call関数を変数に入れ、テスト内で変数にモックを再代入すればOKです。

```go:someprocess.go
package someprocess

func Run(val int) int {
	return call(val)
}

var call = func(val int) int {
	return val
}
```

```go:someprocess_test.go
package someprocess

import (
	"testing"
)

func TestRun(t *testing.T) {
	call = func(val int) int {
		return val + 10
	}

	patterns := []struct {
		val      int
		expected int
	}{
		{2, 12},
		{8, 18},
		{-10, 0},
	}

	for idx, pattern := range patterns {
		actual := Run(pattern.val)
		if pattern.expected != actual {
			t.Errorf("pattern %d: want %d, actual %d", idx, pattern.expected, actual)
		}
	}
}
```
