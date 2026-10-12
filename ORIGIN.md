# コードの移動元

初版の解析・復号処理は[UtauTTS](https://github.com/yh2237/UtauTTS)から分離しました。

| このリポジトリ | 移動元 |
| --- | --- |
| `oto/parse.go` | UtauTTS `911c0fa`の`internal/otofile/parse.go` |
| `oto/parse_test.go` | 同コミットの`internal/otofile/parse_test.go` |
| `textdecode/decode.go` | 同コミットの`internal/oto/decode.go` |

UtauTTSのMIT Licenseと`Copyright (c) 2026 yh`を継承しています。分離にあたり、パッケージ名と公開説明を整理し、単独で実行できる使用例、復号テスト、fuzz試験、CIを追加しました。

UtauTTS固有の`SourceGroup`、音源選択、合成計画、ファイルキャッシュ、音源・モデル・辞書は含みません。

`presamp/parse.go`はUtauTTS `f668c73`時点の`internal/voicebank/presamp.go`から解析処理を分離したものです。分類・置換・語尾の意味は変えず、行・欄・aliasをすべて分割する際の一時配列をなくしました。ファイル探索とfrontendへの変換は含みません。

`prefixmap/parse.go`はUtauTTS `4fb5b9f`時点の`internal/voicebank/metadata.go`からprefix.mapの解析を分離したものです。接辞の空欄・空白、音階名、記載順、不正な行の診断を保ち、行・欄をすべて分割する処理を除きました。最寄り音階の選択は含みません。

`character/parse.go`はUtauTTS `1cec76c`時点の`internal/voicebank/character_yaml.go`と`internal/voicebank/metadata.go`から、character.txtの名前とcharacter.yamlの解析を分離したものです。限定的なYAML表記、音域、行番号付き診断、既定のphonemizerの読込規則を保っています。音階番号への変換は同コミットのmetadata.goを継承しました。

ファイル探索・復号、subbankへのID・パス・順序番号の付与、最寄り音階や色の選択は`character`パッケージに含みません。
