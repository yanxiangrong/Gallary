FROM ubuntu:latest
WORKDIR /go/src/Gallery/
COPY . .
RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN apt-get update && apt-get upgrade -y
RUN DEBIAN_FRONTEND=noninteractive TZ=Asia/Shangshai apt-get install -y tzdata
RUN apt-get install -y ffmpeg libvips-dev golang
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn
RUN go build -o Gallery .
ENTRYPOINT ["./Gallery"]
