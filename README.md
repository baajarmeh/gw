# gw

## Step by step.

dependencies

- mysql
- golint
- redis

Database(MySQL)

``` shell
mysql> create database gwdb;
Query OK, 1 row affected (0.02 sec)

mysql> create user gw@'127.0.0.1' IDENTIFIED BY  'gw@123';
Query OK, 0 rows affected (0.03 sec)

mysql> grant all on gwdb.* to  gw@'127.0.0.1';
Query OK, 0 rows affected (0.00 sec)

mysql> flush privileges;
Query OK, 0 rows affected (0.01 sec)
```


golint

``` shell
go get -u golang.org/x/lint/golint

# references
# https://github.com/golang/lint
```