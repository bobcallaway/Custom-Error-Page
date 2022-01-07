FROM golang:alpine as build-env

RUN apk add git

# Copy source + vendor
COPY . /go/src/github.com/leoquote/custom-error-page
WORKDIR /go/src/github.com/leoquote/custom-error-page

# Build
ENV GOPATH=/go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -v -a -ldflags "-s -w" -o /go/bin/custom-error-page .

FROM scratch
COPY --from=build-env /go/bin/custom-error-page /usr/bin/custom-error-page
ENV ERROR_FILES_PATH=/etc/custom-error-page/templates
ENV GIN_MODE=release
COPY templates ${ERROR_FILES_PATH}
ENTRYPOINT ["custom-error-page"]
