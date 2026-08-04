# Backend

## 開発方法

backendディレクトリで以下を実行すると、Goファイルの変更時にAPIサーバーが再起動します。Task CLIのインストールは不要です。

```sh
go tool air -c .air.toml
```

リポジトリのルートから起動する場合は以下を実行します。

```sh
go -C backend tool air -c .air.toml
```

Task CLIをインストールしている場合は、リポジトリのルートで`task dev-backend`を実行することもできます。

APIサーバーはデフォルトで`http://localhost:8080`に起動します。
