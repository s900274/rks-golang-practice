# magneto

## Requirements

Go Version : >= 1.5

Govendor : https://github.com/kardianos/govendor (Follow the installation guide)

## Get started

1. Clone the repository from the Gitlab to your GOPATH ($GOPAHT/src/gitlab.kingbay-tech.com/engine-lottery/)

2. Move to the repository folder

3. Execute
```
 $ ./scripts/initialization_project.sh <Your_Project_Name>
```

4. Execute govendor sync
```
$ govendor sync
```

5. Generate swagger files
```
$ cd internal/httpservice
```
```
$ ../../tools/swagger/swag init -g httpserver.go -s ../../api/swagger
```

6. Good luck

## Project Layout

### `/api`
放API spec的地方，例如 swagger的docs.json

### `/cmd`
程式的進入點，主要放main跟init，業務邏輯不要放在這裡

### `/config`
設定檔，例如toml, json等等

### `/init`
編譯檔和supervisor放置的地方(就是以前的bin)

### `/internal`
主要業務邏輯放置的地方(不與其他專案共享)

### `/pkg`
共用的包或是我們自己有改過的第三方包(可能與其他專案共享)

### `/scripts`
make, load, initial script

### `/tools`
放置一些tool，例如swag, thrift等等

### `/web`
放置web template或靜態檔案

## Reference
* [Golang Project Layout](https://github.com/golang-standards/project-layout)

* [Govendor](https://github.com/kardianos/govendor)

* [Gin with Swagger](https://github.com/swaggo/gin-swagger)
