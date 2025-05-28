FROM docker.m.daocloud.io/golang:alpine AS builder
# docker buildx build --platform linux/amd64 -t wise:v1.0.0-beta1-0528-01 --load .
# docker tag ea33add688c2 registry.cn-hangzhou.aliyuncs.com/jenkins_construct_images/wise:v1.0.0-beta1-0528-04
# docker push registry.cn-hangzhou.aliyuncs.com/jenkins_construct_images/wise:v1.0.0-beta1-0528-04
# goctl kube deploy --name wise-back --namespace wise --port 8888 --o wise-back-deploy.yaml
LABEL stage=gobuilder

ENV CGO_ENABLED 0


RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o /app/wise wise.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/wise /app/wise
COPY ./etc /app/etc
COPY ./dist /app/dist

CMD ["./wise", "-f", "etc/wise-prod-api.yaml"]
