FROM golang:1.20.5
COPY ./;lemmygousers ./
CMD ["./lemmygousers"]