# get golang container
FROM golang:1.22.5 AS builder

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
FROM alpine:3.20.1 AS app

# create workdir
WORKDIR /opt/app

# add ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# create nonroot user and group
RUN addgroup -S nonroot && \
  adduser -S nonroot -G nonroot && \
  chown -R nonroot:nonroot .

# set user to nonroot
USER nonroot:nonroot

# copy binary from builder
COPY --from=builder --chown=nonroot:nonroot --chmod=544 /go/src/app .

# expose port 8080
EXPOSE 8080

# run application
CMD ["./app"]
