FROM golang:1.24-bookworm
WORKDIR /app 
COPY go.mod go.sum *.go /app
RUN go mod download
ARG CACHE_BUSTER=default
RUN CGOBUILD=0 go build
EXPOSE 8080
CMD ["/app/qxm_backend"]
