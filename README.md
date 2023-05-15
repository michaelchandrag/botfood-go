Installation

```
go mod tidy
source configs/.env
```


Deployment Production
*use Supervisor to run background
```
cd /var/www/botfood-go
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

Nginx Configuration
```
server {
    # listen 80;                # the port nginx is listening on
    server_name     goapi.botfood.id;    # setup your domain here

    gzip            on;
    gzip_types      text/plain application/xml text/css application/javascript;
    gzip_min_length 1000;

    location / {
        expires $expires;

        proxy_redirect                      off;
        proxy_set_header Host               $host;
        proxy_set_header X-Real-IP          $remote_addr;
        proxy_set_header X-Forwarded-For    $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto  $scheme;
        proxy_read_timeout          60m;
        proxy_connect_timeout       60m;
        proxy_pass                          http://127.0.0.1:3001;
    }

    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/api.botfood.id/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/api.botfood.id/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}

server {
    
    if ($host = goapi.botfood.id) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    listen          80;
    server_name     goapi.botfood.id;
    return 404; # managed by Certbot

}


```