# 移動元

初版の解析・復号処理は[UtauTTS](https://github.com/yh2237/UtauTTS)から分離しました。

| このリポジトリ | 移動元 |
| --- | --- |
| `oto/parse.go` | UtauTTS `911c0fa`の`internal/otofile/parse.go` |
| `oto/parse_test.go` | 同コミットの`internal/otofile/parse_test.go` |
| `textdecode/decode.go` | 同コミットの`internal/oto/decode.go` |

UtauTTSのMIT Licenseと`Copyright (c) 2026 yh`を継承しています。パッケージ名と公開説明を整理し、独立した使用例・復号テスト・fuzz試験・CIを追加しています。

UtauTTS固有の`SourceGroup`、音源選択、合成計画、ファイルキャッシュ、音源・モデル・辞書は含みません。
