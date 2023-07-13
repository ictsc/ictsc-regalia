# cspell

<https://cspell.org/>  
スペルチェッカーです。

**以下のルールをよく読んで使って下さい！！！**

## ルール

### 除外ファイルの追加

`cspell.config.yaml`の`ignorePaths`に追加します。

- 追加して良いファイル
  - モジュール管理ファイル
    - `go.mod` `package.json`など
  - 既存リソース / 自動生成されるもの
    - `node_modules` protobufの生成コード など
  - 各種設定ファイル
    - `.gitignore` `.prettierrc` `.editorconfig`など
    - K8sマニュフェストやActions定義などは**含みません**
- 追加してはいけないファイル
  - 所謂コード・ドキュメント類全て
  - 将来のメンテナーが解読する可能性のあるファイル

### 辞書の追加

`dictionary.txt`に追加します。  
VSCodeの場合は、検出されている青線部にマウスをホバーし、`クイックフィックス` -> `Add "～" to dictionary`を選択すると追加できます。

- 追加して良いワード
  - ドメイン用語 (ictscやマイクロサービスのコンポーネントの名前等)
  - 予約語・定義済み識別子
  - モジュール・パッケージの名前 (パスに含まれているもの全部)
    - 自分で定義したパッケージも、最大限の努力の結果ならば許可します
  - モジュール・パッケージに含まれる関数・変数名
  - 既存のツール等の名前
- 追加してはいけないワード
  - 上記以外以外すべて

## 使い方

- GitHub Actionsにより、`main`ブランチへのプッシュ / プルリクエストの作成で自動実行されます
  - 手元で実行はしなくてよいですが、推奨します
- VSCodeの場合
  - 拡張機能から`streetsidesoftware.code-spell-checker`を入れます
    - 推奨拡張機能に入ってます
  - `make init`で初期設定を入れます
  - これで自動でスペルチェックが走ります
- それ以外
  - Nodeを入れます
  - `npm install -g cspell@latest`でcspellを入れます
  - `cspell -c dev/cspell/cspell.config.yaml "**"`でスペルチェックします
