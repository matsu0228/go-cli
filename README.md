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


