# get golang container
FROM golang:1.25.3 AS builder

# get args
ARG apiVersion=unknown

# create and set workingfolder
WORKDIR /go/src/

# copy go mod files and sourcecode
COPY go.mod go.sum ./
COPY src/ .

# download go mods and compile the program
RUN go mod download && \
  CGO_ENABLED=0 GOOS=linux go build \
  -a -installsuffix cgo -ldflags="-w -s \
  -X 'main.apiVersion=${apiVersion}' \
  " -o app ./...


# get alpine container
FROM alpine:3.22.2 AS app

# create workdir
WORKDIR /opt/app

# add packages, create nonroot user and group
RUN apk --no-cache add ca-certificates tzdata && \
  addgroup -S nonroot && \
  adduser -S nonroot -G nonroot && \
  chown -R nonroot:nonroot .

# set user to nonroot
USER nonroot:nonroot

# copy binary from builder
COPY --from=builder --chown=nonroot:nonroot --chmod=555 /go/src/app .

# expose port 8080
EXPOSE 8080

# run application
CMD ["./app"]
