FROM golang:1.18

RUN git -c http.sslVerify=false clone https://github.com/RedHatCRE/syncron.git
WORKDIR /go/syncron/
RUN touch syncron.yaml
RUN go mod tidy
RUN go build cmd/adhoc/main.go

ENTRYPOINT [ "/go/syncron/main" ]
CMD ["download", "insights"]