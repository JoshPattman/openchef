FROM golang:1.24.3-bookworm

RUN mkdir /cgoprecomp
COPY ./cgoprecomp /cgoprecomp
WORKDIR /cgoprecomp
RUN go run .

RUN mkdir /src-code
RUN mkdir /src-code/data
RUN mkdir /src-code/utils

COPY ./servers/data /src-code/data
COPY ./utils /src-code/utils
COPY ./servers/data/datadb /src-code/datadb

WORKDIR /src-code

RUN go work init
RUN go work use data
RUN go work use utils
RUN go work use datadb
WORKDIR /src-code/data
RUN go build -o /data-service-binary .

WORKDIR /

CMD ["/data-service-binary"]