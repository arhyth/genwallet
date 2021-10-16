FROM golang:1.16
ADD . /workspace

# build executable files
WORKDIR /workspace
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /release/genwallet . \
    && chmod +x /release/genwallet

EXPOSE 8000
ENTRYPOINT [ "/release/genwallet" ]
