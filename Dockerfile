FROM golang as builder
WORKDIR /build/
COPY . /build/
RUN go build -o /build/bin/tictactoe .

FROM alpine:latest
WORKDIR /app/
RUN apk add libc6-compat
COPY --from=builder /build/bin/tictactoe /app/tictactoe
EXPOSE 8888
ENTRYPOINT [ "/app/tictactoe" , "server"]