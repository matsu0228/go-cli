#!/bin/bash

# conf

# move
dd if=/dev/zero of=./log_big_old.log bs=1K count=10
touch -t 201712201000 ./log_big_old.log

# dont move
dd if=/dev/zero of=./big.log bs=1K count=10
touch -t 201712201000 ./old.log
touch ./invalid_extenstion.txt
