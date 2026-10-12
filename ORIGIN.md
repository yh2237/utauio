# 移動元

初版の解析・復号処理は[UtauTTS](https://github.com/yh2237/UtauTTS)から分離しました。

| このリポジトリ | 移動元 |
| --- | --- |
| `oto/parse.go` | UtauTTS `911c0fa`の`internal/otofile/parse.go` |
| `oto/parse_test.go` | 同コミットの`internal/otofile/parse_test.go` |
| `textdecode/decode.go` | 同コミットの`internal/oto/decode.go` |

UtauTTSのMIT Licenseと`Copyright (c) 2026 yh`を継承しています。パッケージ名と公開説明を整理し、独立した使用例・復号テスト・fuzz試験・CIを追加しています。

UtauTTS固有の`SourceGroup`、音源選択、合成計画、ファイルキャッシュ、音源・モデル・辞書は含みません。

`presamp/parse.go`はUtauTTS `f668c73`時点の`internal/voicebank/presamp.go`の解析処理から分離しました。分類・置換・語尾の意味を維持し、行・欄・aliasの全分割による一時配列をなくしています。ファイル探索とfrontendへの変換は含みません。

`prefixmap/parse.go`はUtauTTS `4fb5b9f`時点の`internal/voicebank/metadata.go`のprefix.map解析から分離しました。接辞の空欄・空白、音階名、記載順・不正行診断を保持し、行・欄の全分割を除いています。最寄り音階の選択は含みません。

`character/parse.go`はUtauTTS `1cec76c`時点の`internal/voicebank/character_yaml.go`と`internal/voicebank/metadata.go`から、character.txtの名前とcharacter.yamlの解析を分離しました。限定的なYAML表記、音域、行診断、既定phonemizerの読込規則を保持しています。音階番号への変換は同コミットのmetadata.goを継承しました。ファイル探索・復号・subbankのID/パス/順序番号の付与・最寄り音階や色の選択は含みません。
