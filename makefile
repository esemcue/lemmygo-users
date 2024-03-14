.DEFAULT_GOAL := deploy

deploy: 
	go build -o ./lemmygousers
	docker build -t registry.gitlab.com/lemmygo/lemmygo-users .
	docker push registry.gitlab.com/lemmygo/lemmygo-users
	rm ./lemmygousers