default: deploy

deploy:
	docker-compose up -d --force-recreate --build

start: build
	docker-compose up -d
stop:
	docker-compose down

build:
	docker-compose build --force-rm

debug: start
	sleep 5
	docker logs detect-vpn

logs:
	docker logs detect-vpn

clean:
	docker system prune -f
	-rm -f detect-vpn/detect-vpn
	-rm -f publisher/publisher

test:
	docker-compose up -d mosquitto
	go test -timeout 2m -run ^Test*$ github.com/jxsl13/tw-moderation/common/mqtt

update:
	# first delete the required dependencies and then execute this
	go get -u

proxy:
	# create a temporary proxy that exposes the rabbitmq broker on port 5673
	socat tcp-listen:5673,reuseaddr,fork tcp:localhost:5672


