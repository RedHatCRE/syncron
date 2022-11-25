FROM quay.io/projectquay/golang:1.18

#RUN git -c http.sslVerify=false clone https://github.com/RedHatCRE/syncron.git
#COPY credentials /.aws/credentials
WORKDIR /go/syncron/
RUN touch syncron.yaml
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1
RUN go install github.com/redhatcre/syncron@latest
#ENTRYPOINT [ "syncron" ]
#CMD ["download", "insights"]