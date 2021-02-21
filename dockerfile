FROM golang:1.15.8-alpine AS build
WORKDIR /src
COPY . .
RUN export CGO_ENABLED=0
RUN go build -o server ./app
FROM alpine:latest AS bin
WORKDIR /root/
COPY --from=build /src/server .
CMD ["./server"]

