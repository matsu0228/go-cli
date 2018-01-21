# log dir

- this is directory your dumy log storage
- please create dumy log file

```
dd if=/dev/zero of=log/log_big_old.log bs=1K count=10

touch -t 201712201000 log/old.log

touch log/invalid_extenstion.txt
```
