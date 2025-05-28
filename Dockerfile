FROM docker.m.daocloud.io/golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0


RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o /app/wise wise.go


FROM docker.m.daocloud.io/golang:alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/wise /app/wise
COPY ./etc /app/etc
COPY ./dist /app/dist

CMD ["./wise", "-f", "etc/wise-prod-api.yaml"]
