# CLI-Chat-Server
- CLIで動作するチャットアプリのサーバープロジェクト
- [CLI-Chat-Client](https://github.com/ChikaKakazu/CLI-Chat-Client)も動作する必要がある

## 起動
1. dockerコンテナを起動する
    ```
    docker-compose up -d
    ```
1. Client側のdockerコンテナを起動する
2. コンテナ内で実行コマンドを叩く
   `go run main.go`