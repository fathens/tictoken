# tictoken

Solidity のソースをコンパイルして ethereum にデプロイするツール。

## Usage

main.go を起点とする。

```shell
go run main.go [options] <command> [args ...]
```

与えられるオプションは以下の通り

| オプション名 | 説明                      | デフォルト値     |
| ------------ | ------------------------- | --------------   |
| --config     | 設定ファイルのパス        | .config.toml     |
| --hdpath     | マスターキーからのHDパス  | m/44'/60'/0'/0/0 |
| --solc       | Solidity コンパイラのパス | solc             |



#### 環境変数

| 環境変数名        | 説明                          |
| ----------------- | ----------------------------- |
| TICTOKEN_MNEMONIC | マスターキーのための Mnemonic |



#### 設定ファイル

```
rpcserver = "http://12.34.56.78:8545"
privatekey = "<hex of ECDSA secp256k1>"
```

`privatekey` が存在する場合は Mnemonic よりも優先される。


## コマンド

### deploy

```shell
go run main.go [options] deploy <solファイルのパス> [コンストラクタへの引数 ...]
```

成功すれば、デプロイされたコントラクトアドレスを出力する。

### invoke

```shell
go run main.go [options] invoke <contract address> <ABIファイルのパス> <メソッド名> [メソッドへの引数 ...]
```

* ABI ファイルは JSON 形式であること。
* 引数の型は今のところ `string` と `address` のみをサポートしている。

成功すれば、戻り値を出力する。
