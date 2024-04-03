FROM golang:alpine as builder

LABEL maintainer="Yigithan Karabulut <yigithannkarabulutt@gmail.com>"

RUN apk update && apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
COPY s3-objects.yml ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/s3-objects.yml .


ENV AWS_REGION=${AWS_REGION}
ENV AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
ENV AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}

ENV DB_NAME=${DB_NAME}
ENV DB_USER=${DB_USER}
ENV DB_PASS=${DB_PASS}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_PRODUCT_COLLECTION=${DB_PRODUCT_COLLECTION}
ENV DB_OBJECTINFO_COLLECTION=${DB_OBJECTINFO_COLLECTION}

CMD ["./main"]