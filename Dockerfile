# get golang container
FROM golang:1.22.1 as builder

# get args
ARG apiVersion=unknown

# create and set workingfolder
WORKDIR /go/src/

# copy go mod files
COPY go.mod go.sum ./

# download go mods
RUN go mod download

# copy all sourcecode
COPY src/ .

# compile the program
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'main.apiVersion=${apiVersion}'" -o app ./...


# get alpine container
FROM alpine:3.19.1 as app

# create nonroot user and group
RUN addgroup -S nonroot \
  && adduser -S nonroot -G nonroot

# add ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# create workdir
WORKDIR /root/

# copy binary from builder
COPY --from=builder /go/src/app .

# set user
USER nonroot

# expose port 8080
EXPOSE 8080

# run application
CMD ["./app"]
