# utauio

UTAU音源のファイルを扱う小さなGoライブラリです。初版は`oto.ini`の解析と日本語テキストの復号を提供します。Go 1.25以上で利用できます。

UtauTTSから分離した処理を、合成エンジン・辞書・GUI・音源データに依存せず利用できます。

## インストール

```sh
go get github.com/yh2237/utauio@v0.1.0
```

| パッケージ | 用途 | 依存 |
| --- | --- | --- |
| `oto` | 復号済みテキストの解析 | 標準ライブラリのみ |
| `textdecode` | UTF-8/BOM、UTF-16LE/BE、Shift_JISの復号 | `golang.org/x/text` |

## 使用例

```go
data, err := os.ReadFile("voice/oto.ini")
if err != nil {
    return err
}
text, encoding, err := textdecode.Decode(data)
if err != nil {
    return err
}
diagnostics, err := oto.Scan(strings.NewReader(text), "voice", func(entry oto.Entry) {
    fmt.Printf("%s: %s (line %d)\n", entry.Alias, entry.Filename, entry.Line)
})
if err != nil {
    return err
}
fmt.Println(encoding, diagnostics)
```

importは`github.com/yh2237/utauio/oto`と`github.com/yh2237/utauio/textdecode`です。実行可能な例は`go run ./examples/oto /path/to/oto.ini`で試せます。

## 解析の契約

`oto.Scan(reader, baseDir, visit)`は復号済みのテキストを順番に読みます。

- 正常行を`visit`へ渡します。同じaliasの行もすべて渡し、順序を保ちます。
- 空行と`#`・`;`で始まるコメント行を読み飛ばします。
- aliasが空なら録音名から補完し、空の数値は0とします。
- 5パラメータはミリ秒です。負のblank等の形式上の値を保ち、NaN/Infは診断として拒否します。
- 不正行は行番号と理由を返し、後続行の解析を続けます。readerの失敗と約1MiBの行長上限はエラーです。
- `visit`は必須です。エラーが返るまでに正常行のコールバックが実行済みの場合があります。
- 相対録音パスは`baseDir`基準で解決します。基準を省略すると相対のままで、作業ディレクトリを暗黙に使いません。
- パス解決は実行OSの`filepath`規則です。相対パスのWindows区切りを補完しますが、他OSのドライブ名の解釈は保証しません。
- ファイルの存在や音源ルート内かどうかは検査しません。必要な検査は利用側で行います。
- aliasと5パラメータより後の欄は無視します。書き戻しや元のコメント・表記の保持には対応していません。

## 復号の契約

`textdecode.Decode`は復号結果と文字コード名を返します。BOM付きUTF-16、UTF-8（BOMは除去）、それ以外をShift_JISとして扱います。文字コード名は`UTF-8`、`UTF-16LE`、`UTF-16BE`、`Shift_JIS`です。

任意の文字コードを推定するAPIではありません。UTF-16の奇数長を拒否し、未対応の符号列の扱いはGoのUTF-16変換と`x/text`の復号器に従います。UtauTTSの従来動作を維持しています。

## 検証

```sh
go test ./...
go vet ./...
go test ./oto -run '^$' -fuzz FuzzScan -fuzztime 10s
```

CIにはWindows・Linux・macOS、Go 1.25とstable、Go wasmビルドの検査を設定しています。公開バージョンは`v0.1.0`です。

## 出典とライセンス

MIT Licenseです。元の著作権表示を継承しています。移動元と対応ファイルは[ORIGIN.md](ORIGIN.md)、依存の通知は[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)に記載しています。
