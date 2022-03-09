# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.17-alpine

# Set working directory
WORKDIR /app

# Run this with docker build --build_arg $(go env GOPROXY) to override the goproxy
# ARG goproxy=https://proxy.golang.org
# ENV GOPROXY=$goproxy

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
go mod download

# Copy all sources
COPY ./ ./

RUN go build -o /main

EXPOSE 8081

CMD ["ls"]
CMD ["tree"]
CMD [ "/main", "-host", "0.0.0.0" ]