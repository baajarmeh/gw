# gw

gw (golang web), gin + gorm + go-redis + mysql

out-of-the-box, Fully features:

- Auth manager
- Permission manager
- Multi-tenancy
- ORM (gorm. mysql, postgresql)
- Cache (redis)
- Modular module application.
- REST/Dynamic style Api support.
- Out-of-the-box Web Api framework.
- Api decorator.
- Base on go template config generator.

## Step by step.

dependencies

- mysql
- golint
- redis
- openssl certificates

Database(MySQL)

``` shell
brew install mysql
brew services start mysql

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

redis

``` shell
brew install redis
brew services start redis
```

certs

``` shell
mkdir -p config/etc
openssl genrsa -out config/etc/gw.key 2048
openssl rsa -in config/etc/gw.key -pubout -out config/etc/gw.pem
```

## Quick Start

```shell script
wget -O $GOROOT/bin/gwcli https://raw.githubusercontent.com/oceanho/gw/dev/cmd/gwcli/scripts/cli.sh
chmod +x $GOROOT/bin/gwcli
```

### Create Project

```shell script
mkdir -p ~/workspace/myapp
gwcli newproject mywebapi .
```
