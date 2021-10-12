ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine AS build

ENV GO111MODULE=on \
    APP_BUILD_PATH="/app" \
    APP_BUILD_NAME=creativeadvtech \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR ${APP_BUILD_PATH}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH}  go build  -o ${APP_BUILD_NAME}
RUN chmod +x ${APP_BUILD_NAME}


FROM alpine

LABEL Name="creativeadvtech" \
      Version=0.0.1

ENV NAME=app \
    APP_BUILD_PATH="/app" \
    APP_BUILD_NAME=creativeadvtech

COPY entrypoint.sh /usr/local/bin/
RUN  chmod +x /usr/local/bin/entrypoint.sh

COPY --from=build ${APP_BUILD_PATH}/${APP_BUILD_NAME} /usr/local/bin/

ENTRYPOINT ["entrypoint.sh"]