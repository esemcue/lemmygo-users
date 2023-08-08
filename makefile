.DEFAULT_GOAL := deploy

deploy: 
	go build -o ./lemmygousers
	docker build -t registry.gitlab.com/lemmygo/lemmygousers .
	docker push registry.gitlab.com/lemmygo/lemmygousers
	rm ./lemmygousers