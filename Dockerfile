FROM golang:1.13

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# go mod need git
RUN mkdir /project
WORKDIR /project

# 先下载依赖再拷贝代码，可以有效利用缓存
COPY go.mod .
COPY go.sum .

# download 3rd package
RUN export GOPROXY=https://goproxy.cn && go mod download

# copy project
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix -o project

# run
ENTRYPOINT ./project