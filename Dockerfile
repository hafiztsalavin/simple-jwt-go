############################
# build executable binary
############################

FROM golang:1.17-alpine 

RUN apk update && apk upgrade && \
    apk --update add git make && apk add --no-cache git

WORKDIR /app

COPY . .

RUN make clean

RUN make engine

CMD ["./engine"]
