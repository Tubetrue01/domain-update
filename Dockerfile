FROM alpine:3.12.1
ADD ./proc /usr/local/app/proc

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

ENTRYPOINT ["sh", "-c","./usr/local/app/proc -e xxxx@qq.com:xxxx@qq.com -a xxxxx:xxxxx"]

# docker build -t proc . --network=host
# docker run -d --restart=always --net=host --name=proc proc
