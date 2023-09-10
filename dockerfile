FROM golang:1.20.5
COPY ./lemmygousers ./
EXPOSE 8081
CMD ["./lemmygousers"]