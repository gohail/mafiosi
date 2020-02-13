# Mafiosi

### How to start server

```bash
go build -o app  
./app
```

also
```bash
cd bin
./app
```
Use -host flag if needed 
```
./app -host=yourhost:port 
```
default:"localhost:8080"


### Connection to WS


### Docker container
to create container of go application
```bash
docker build . -t mafiosi:go-back
```
---
to start docker container
```bash
docker run -p 8080:8080 mafiosi:go-back
```
---
more information about installs and settings docker by [link](https://docs.docker.com/install/linux/docker-ce/ubuntu/)