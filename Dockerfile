From golang:1.14.6-buster as builder

WORKDIR /go/src/stella
COPY . /go/src/stella

RUN go get -d -v
RUN go build -o /go/bin/stella


FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/bin/stella /
CMD ["/stella"]
