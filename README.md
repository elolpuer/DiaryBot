<h1>Diary Bot</h1>

<h2><a href="https://t.me/godofdiarybot">Link</a></h2>

<h2>Technology stack</h2>

Platform: <a href="https://core.telegram.org/bots" target="_blank">Telegram Bots</a><br>
Language: <a href="https://go.dev/" target="_blank">Go</a><br>
Database: <a href="https://www.postgresql.org/" target="_blank">PostgreSQL</a><br>
In-memory data store: <a href="https://redis.io/" target="_blank">Redis</a>

<h2>Manual deployment</h2>

<h3>Connect to the server and it's configuration</h3>

```shell
ssh root@your_server_ip
adduser egorg
usermod -aG sudo egorg
```

<h3>Installing packages and creating a database</h3> 

<h4>Postgres</h4> 

```shell
sudo apt update 
sudo apt install postgresql postgresql-contrib
sudo -i -u postgres
psql 
create database somedb;
createuser egorg;
ALTER USER egorg WITH ENCRYPTED PASSWORD 'pass';
alter user egorg superuser createrole somedb;
alter database somedb owner to egorg;
\q
psql -h localhost -p 5432 -U egorg somedb
\i pathToBot/cfg/db.sql
```


<h4>Golang</h3>

```shell
sudo apt install golang
go version
mkdir ~/workplace/src/github.com/elolpuer
nano ~/.profile
```

```shell
export GOPATH=$HOME/workplace\
export PATH=$PATH:$GOPATH/bin\
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin
~/.profile
```
<h4>Redis</h4>

```shell
sudo apt install redis-server
supervised systemd
sudo systemctl restart redis.service
sudo systemctl status redis
```

<h4>Git</h4> 

```shell
sudo apt install git
``` 

<h3>Clonning project from github and building</h3> 

```shell
cd ~/workplace/src/github.com/elolpuer
git clone https://github.com/elolpuer/DiaryBot.git
go build main.go
``` 


<h3>Running daemon</h3> 

```shell 
cd ../../etc/systemd/system/
sudo nano gosomething.service
``` 
Insert this:
```shell
[Unit]\
Description = Something Description

[Service]\
WorkingDirectory=/root/workplace/src/github.com/elolpuer/someproject\
ExecStart=/root/workplace/src/github.com/elolpuer/someproject/main\
User=root\
Group=root\
Restart=always

[Install]\
WantedBy=multi-user.target
```

<h3>Reboot and run daemon</h3>

```shell
sudo reboot
ssh root@your_server_ip
cd ../../etc/systemd/system/
sudo systemctl enable gosomething.service
sudo systemctl daemon-reload
sudo systemctl start gosomethint.service
```

