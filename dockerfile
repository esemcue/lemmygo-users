FROM golang:1.20.5
COPY ./lemmygousers ./
COPY ./.env ./deploy.env
CMD ["./lemmygousers"]