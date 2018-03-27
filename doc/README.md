# Go言語導入

## Go言語学習にあたって

* go-tutorial(公式)
* チートシート: https://github.com/a8m/go-lang-cheat-sheet
* goによるweb開発の基本(slide): https://www.slideshare.net/hico_horiuchi/goweb-58764875
* 巨大なライブラリ集: https://github.com/avelino/awesome-go

## Goの開発環境

* コーディング規約
  * gofmt(エディタで保存時に自動で実行されるようにしておく or push時に自動テスト)
  * golang CodeReviewComments(日本語訳): https://qiita.com/knsh14/items/8b73b31822c109d4c497
* package管理
  * godep
* DB周り
  * ORM: GORM(msssqlにも対応)
  * Migration: gooseがよさげ。GORMでもmigrationがあるので、不要？
* web application
  * WAF(Framework)は利用しない。下記を組み合わせて実装する
  * net/http: http通信のためのライブラリ
  * html/temlate or Ace: templateエンジン
* logging: hashicorp/logutils

### godep

* install & setup
```
# go get でインストール
go get -u github.com/golang/dep/cmd/dep

# gopath配下のプロジェクトでないと、dep initできない
go env
```

* 使い方
```
# 初期設定
dep init

# パッケージをvenderディレクトリ配下に作成
dep ensure

# 新しいパッケージを依存関係に追加しインストールを行う
dep ensure -add github.com/foo/bar

# パッケージをアップデートする
dep ensure -update github.com/foo/bar

# 現在の状態を出力する
dep status
```