FROM golang:1.17
WORKDIR /go/src/Gallery/
COPY . .
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn
RUN apt update && apt upgrade -y && apt install -y ffmpeg libvips-dev
RUN go build -o Gallery .
ENTRYPOINT ["./Gallery"]
