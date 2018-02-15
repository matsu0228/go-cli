# install

## gvm

- `gvm` is version managemt tool: https://github.com/moovweb/gvm
```
## see latest method in upper link ##

[install]
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source **/gvm/**

[check]
which gvm
gvm version
```
- install latest go via gvm (ref: https://gist.github.com/d2s/6503f815431d1587c28bc37bfd715dbf)
```
gvm listall
# go1.4 is requred for latest version
gvm install go1.4
gvm install go1.*.*

[check]
gvm list
gvm use go1.*.* -default
```
- set package env
```
mkdir -p $HOME/hoge/{pkg,bin,src}
gvm pkgset create hoge

gvm pkgenv hoge
# line12
export GOPATH; GOPATH="/Users/nghinv/.gvm/pkgsets/go1.7/hoge:$HOME/hoge:$GOPATH"
# line16
export PATH; PATH="/Users/nghinv/.gvm/pkgsets/go1.7/hoge/bin:${GVM_OVERLAY_PREFIX}/bin:$HOME/hoge/bin:${PATH}"

gvm pkgset use hoge
```

## glide

-package management tool
- https://github.com/Masterminds/glide
- warning: you should edit install shell about $GOPATH
```
curl https://glide.sh/get > glide.sh
vim glide_install.sh

GOPATH='~/hoge'
$GOPATH -> ${GOPATH}

sh glide_install.sh
```

- create
```
cd /your/working/dir/
glide create

cat glide.yaml
```

## revel

- install
```
$ go get github.com/revel/cmd/revel
```
- setup $ run
```
$ revel new github.com/*user*/*myproj*
$ revel run github.com/*user*/*myproj*
```

- glide.yaml
```
$ cd src/github.com/*user*/*myproj*
$ glide create
$ glide install
```

## GORM

- install
```
go get -u github.com/jinzhu/gorm
```

- create db conf
```
$ vim conf/db.fonf

[database]
db.host:       0.0.0.0
db.port:       3306
db.user:       your_user
db.password:   ***
db.name:       your_db
db.protocol:   tcp
db.args:       charset=utf8&parseTime=True&loc=Local
```