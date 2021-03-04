package pkg

// 構造体ではないのでメソッド生成対象外
const a = 100

// 構造体ではないのでメソッド生成対象外
var b int = 200

// 構造体ではないのでメソッド生成対象外
type c int

// typeで型定義したstructのみ生成対象にしているため、この無名構造体は対象外
var d = struct{ A int }{}

// PrintType 自動生成したメソッドを呼び出す
func PrintType() {
	StructA{}.PrintType()
	structB{}.PrintType()
}
