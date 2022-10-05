FROM golang:1.18

RUN git -c http.sslVerify=false clone https://github.com/RedHatCRE/syncron.git
WORKDIR /go/syncron/
RUN touch syncron.yaml
RUN go mod tidy
RUN go install cmd/adhoc/main.go
