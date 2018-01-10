## ftp

- setup ftp-server
```
yum install -y vsftpd
systemctl start vsftpd.service
systemctl enable vsftpd.service

[check]
rpm -qa | grep vsftpd
systemctl status vsftpd.service
systemctl list-unit-files -t service | grep vsftpd
```
- setup ftp-client
```
yum install -y ftp
```

- test
```
ftp localhost

[try when error]
firewalld / Selinux
```
