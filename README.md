# learning-angular-container

## nginx container

Dockerfile
```
FROM nginx
COPY dist/learning-angular /usr/share/nginx/html
```

shell
```
docker build -t content-nginx .
docker run --name content-nginx -d -p 9800:80 content-nginx
```

browser
```
http://localhost:9800/demo.html
http://localhost:9800/index.html
```

## golang container

Dockerfile
```

```

shell
```
docker build -t content-golang .
docker run --name content-golang -d -p 9800:80 content-golang
```

browser
```
http://localhost:9800/demo.html
http://localhost:9800/index.html
```
