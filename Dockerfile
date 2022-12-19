FROM golang:1.20-rc-bullseye

COPY ./build/ ./

CMD ["./main"]
