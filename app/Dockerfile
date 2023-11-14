FROM golang:1.21 as build_base

WORKDIR /tmp/chat

COPY ./* ./

RUN CGO_ENABLED=0 go build -tags netgo,osusergo,static_build -o /chat .

FROM alpine:3.18.4

COPY --from=build_base /chat /chat

ENV SOUTHPARKCHAT_ADDR "0.0.0.0:8888"
ENV SOUTHPARKCHAT_DB_ADDR ":memory:"

EXPOSE 8888/tcp

CMD ["/chat"]
