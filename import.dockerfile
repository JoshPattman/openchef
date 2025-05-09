FROM golang:1.24.3-bookworm

RUN mkdir /cgoprecomp
COPY ./cgoprecomp /cgoprecomp
WORKDIR /cgoprecomp
RUN go run .

RUN mkdir /src-code
RUN mkdir /src-code/import
RUN mkdir /src-code/utils

COPY ./servers/import /src-code/import
COPY ./utils /src-code/utils

WORKDIR /src-code

RUN go work init
RUN go work use import
RUN go work use utils
WORKDIR /src-code/import
RUN go build -o /import-service-binary .

WORKDIR /

CMD ["/import-service-binary"]