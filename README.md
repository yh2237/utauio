# utauio

UTAU音源の設定ファイルを解析するGoライブラリです。`oto.ini`・`presamp.ini`・`prefix.map`・`character.txt`・`character.yaml`の解析と、日本語テキストの復号に対応しています。Go 1.25以上で利用できます。

UtauTTSから分離したライブラリで、合成エンジン・辞書・GUI・音源データには依存しません。

## インストール

```sh
go get github.com/yh2237/utauio@v0.4.0
```

| パッケージ | 用途 | 依存 |
| --- | --- | --- |
| `oto` | 復号済みの`oto.ini`の解析 | 標準ライブラリのみ |
| `presamp` | 母音・子音分類、置換、語尾設定の解析 | 標準ライブラリのみ |
| `prefixmap` | 音階名と接頭辞・接尾辞の解析 | 標準ライブラリのみ |
| `character` | 名前、既定のphonemizer、subbanksの限定的な解析 | 標準ライブラリのみ |
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

この例では`github.com/yh2237/utauio/oto`と`github.com/yh2237/utauio/textdecode`をimportします。実行できる例は`go run ./examples/oto /path/to/oto.ini`で試せます。

## oto.iniの解析

`oto.Scan(reader, baseDir, visit)`は復号済みのテキストを順番に読みます。

- 正常な行を記載順に`visit`へ渡します。同じaliasの行もすべて渡します。
- 空行と`#`・`;`で始まるコメント行を読み飛ばします。
- aliasが空なら録音名から補完し、空の数値は0とします。
- 5つの数値パラメータの単位はミリ秒です。負のblankなど、形式上有効な値は保持します。NaN/Infは診断を返して拒否します。
- 不正な行は行番号と理由を診断として返し、後続の行を読み続けます。readerの読込失敗と、約1MiBの行長上限を超える入力はエラーです。
- `visit`は必須です。エラーを返した時点で、正常な行のコールバックがすでに実行されている場合があります。
- 相対録音パスは`baseDir`を基準に解決します。`baseDir`を省略すると相対パスのまま返し、作業ディレクトリは暗黙に使いません。
- パスの解決には実行OSの`filepath`の規則を使います。相対パスのWindows形式の区切りを補完しますが、他OSでのドライブ名の解釈は保証しません。
- ファイルが存在するか、音源ルート内にあるかは検査しません。必要な検査は呼出し側で行います。
- aliasと5つの数値パラメータより後の欄は無視します。書き戻しや、元のコメント・表記の保存には対応していません。

## 文字コードの復号

`textdecode.Decode`は復号した文字列と文字コード名を返します。BOM付きのUTF-16、UTF-8の順に判定し、それ以外はShift_JISとして扱います。UTF-8のBOMは除去します。返す文字コード名は`UTF-8`、`UTF-16LE`、`UTF-16BE`、`Shift_JIS`です。

任意の文字コードを推定するAPIではありません。UTF-16のデータ長が奇数の場合は拒否します。未対応の符号列の扱いはGoのUTF-16変換と`x/text`の復号器に従い、UtauTTSの従来動作を保っています。

## presamp.ini

`presamp.Parse(text)`は復号済みの文字列から`Config`を返します。

- `[VOWEL]`、`[CONSONANT]`、`[REPLACE]`、`[ENDTYPE...]`、`[ENDFLAG]`を解析します。セクション名の大文字小文字は区別しません。
- aliasをキー、分類名を値とするmapを返します。aliasが重複すると後の行を採用し、空のaliasは除きます。
- 母音分類は`class=reading=aliases`、子音分類は`class=aliases`の欄を使用し、後続欄を無視します。aliasはカンマ区切りです。
- REPLACEは最初の`=`で区切り、キーと値の空白を除きます。
- 語尾はENDTYPEの番号順に返し、同じ番号なら記載順を保ちます。ENDFLAGで対応するビットが指定された語尾を採用します。ENDFLAGが0・省略・不正数値の場合は全語尾を採用し、ENDTYPEの番号が不正または省略された場合は1として扱います。
- 空行、`#`/`;`で始まる行、不正な分類行、未知のセクションは無視します。行末コメントの解釈や診断の生成は行いません。
- 元のテキストやコメントの保存、書き戻し、ファイルの読込、発音生成には対応していません。入力の復号には`textdecode`などを使えます。

実行例は`go run ./examples/presamp`です。

## prefix.map

`prefixmap.Scan(text, visit)`は復号済みの文字列から正常行を記載順に渡し、行番号付き診断を返します。

- タブ区切りの音階名・接頭辞・接尾辞を使用します。音階名は周囲の空白を除いて大文字化します。
- 接辞の空欄と空白を保持します。接尾辞欄は省略可能で、4欄目以降は無視します。
- 音階名が重複する行もすべて渡します。どの行を採用するかは呼出し側で決めます。
- 空行と、`#`/`;`で始まるコメント行は無視します。タブのない不正な行は診断を返し、後続の行を読み続けます。
- 音階名が音楽的に有効か、ファイルが存在するかは検査せず、原音の選択も行いません。空の音階名も既存形式の値として渡します。
- `visit`は必須です。ファイル探索・復号・書き戻しはこのパッケージに含みません。

実行例は`go run ./examples/prefixmap`です。

## character.txt・character.yaml

`character`パッケージは`v0.4.0`から利用できます。

`character.ParseText(text)`は名前を返します。`name`の大文字小文字を区別せず、周囲の空白を除きます。重複時は最後の非空値を採用します。名前以外の項目は無視します。

`character.ParseYAML(text)`は`YAML`と行番号付き診断を返します。

- UtauTTSから継承した限定的な表記を扱います。完全なYAMLパーサーではなく、アンカー・複数行文字列・任意のmappingには対応しません。
- `subbanks:`以降の`- color:`で項目を開始し、`color`・`prefix`・`suffix`・`tone_ranges`を記載順に読みます。これらのキーは大文字小文字を区別します。
- 音域は`C3-C5`などの両端を含む範囲です。ブロック形式のリストと、`[C3-C5, Db4-F#4]`形式を扱い、MIDI音階番号で返します。不正な音域や逆順の音域は診断を返し、後続を読み続けます。
- 有効な音域がない項目は`Low=-1073741824`、`High=1073741823`とします。これは従来の全音域指定を保つ値です。
- subbankのスカラーでは、Go形式の二重引用符による表記と、単一引用符内の`''`を扱います。引用符の外で空白の後に置かれた`#`はコメントとして除きます。
- `default_phonemizer`は大文字小文字を区別せず、最初の値を採用します。従来の読込規則に合わせ、位置やインデントを検査せず、行末コメントも除きません。
- 未知の項目は無視します。ファイル探索、復号、ID・パスの付与、音源選択、書き戻しは含みません。

実行例は`go run ./examples/character`です。

## 検証

```sh
go test ./...
go vet ./...
go test ./oto -run '^$' -fuzz FuzzScan -fuzztime 10s
go test ./presamp -run '^$' -fuzz FuzzParse -fuzztime 10s
go test ./prefixmap -run '^$' -fuzz FuzzScan -fuzztime 10s
go test ./character -run '^$' -fuzz FuzzParse -fuzztime 10s
```

CIにはWindows・Linux・macOS、Go 1.25とstable、Go wasmビルドの検査を設定しています。公開バージョンは`v0.4.0`です。

## 出典とライセンス

MIT Licenseで公開し、元の著作権表示を継承しています。移動元と対応ファイルは[ORIGIN.md](ORIGIN.md)、依存ライブラリのライセンスは[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)に記載しています。
