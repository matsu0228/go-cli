# example of go-cli

- go cli with a library of `spf13/cobra`

# how to build

```
$ pwd
-> (your go-cli dir)

$ go install go-cli
$ go-cli

ls $GOPATH\bin
```

# usage

```
$ go-cli -h
-> check all command

$ go-cli parsertest hoge
$ go-cli parsertest hoge -t 5
-> example of argment parser
```


# go test

- sandbox/example
- create *target_file*_test.go
    - `import "testing"`
- exec go test
```
go test packagePath  # cf. go test github.com/hoge/moge パッケージを指定
go test -run ''      # カレントディレクトリ内の全てのテストを実行
go test -run Foo     # Fooを含むテスト関数を実行
```
- ref: https://qiita.com/taizo/items/82930518430f940721a0
- tableDrivenTests: https://github.com/golang/go/wiki/TableDrivenTests
