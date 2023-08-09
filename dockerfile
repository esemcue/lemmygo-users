FROM golang:1.20.5
COPY ./lemmygousers ./
COPY ./config.yaml ./
CMD ["./lemmygousers"]