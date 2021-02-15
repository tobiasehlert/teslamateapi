# get latest 1.15.8 golang
FROM golang:1.15.8

# create and set workingfolder
WORKDIR /go/src/

# download dependencies
RUN go get -v -u -d github.com/gin-gonic/gin github.com/lib/pq github.com/eclipse/paho.mqtt.golang github.com/go-sql-driver/mysql

# copy all sourcecode
COPY src/ .

# compile the program
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o TeslaMateApi .


# get latest alpine container
FROM alpine:latest

# add ca-certificates
RUN apk --no-cache add ca-certificates

# create workdir
WORKDIR /root/

# copy binary from first container
COPY --from=0 /go/src/TeslaMateApi .

# run application
CMD ["./TeslaMateApi"]
