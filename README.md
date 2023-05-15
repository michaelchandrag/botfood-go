Installation

```
go mod tidy
source configs/.env
```


Deployment Production
*use Supervisor to run background
```
cd /var/www/botfood-go
source .env
GOOS=linux GOARCH=amd64 go build -o build cmd/http-botfood/main.go
```

Supervisor configuration
```
[program:botfood-go]
directory=/var/www/botfood-go
command=/var/www/botfood-go/build/main
autostart=true
autorestart=true
stderr_logfile=/var/log/botfood-go/api.err
stdout_logfile=/var/log/botfood-go/api.log
startretries=3
```