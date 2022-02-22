FROM golang:1.17
WORKDIR /go/src/Gallery/
COPY . .
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn
RUN go build -o Gallery .
RUN apt update && apt install -y ffmpeg
ENTRYPOINT ["./Gallery"]
