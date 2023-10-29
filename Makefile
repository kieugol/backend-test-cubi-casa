build-app:
	docker-compose build --no-cache app
start:
	docker-compose up -d
restart:
	docker restart rest-api
stop:
	docker-compose down
logs:
	docker logs -f rest-api
ssh-app:
	docker exec -it rest-api bash
test-hub:
	go test -v ./tests/unit/hub.service_test.go \
  		-coverpkg=./api/v1/services,./api/v1/repository \
  		-coverprofile=./tests/report/hub/hub.service.coverage.out
	go tool cover -html ./tests/report/hub/hub.coverage.out \
		-o ./tests/report/hub/hub.service.coverage.html
test-team:
	go test -v ./tests/unit/team.service_test.go \
  		-coverpkg=./api/v1/services,./api/v1/repository \
  		-coverprofile=./tests/report/team/team.service.coverage.out
	go tool cover -html ./tests/report/team/team.service.coverage.out \
		-o ./tests/report/team/team.service.service.coverage.html
test-user:
	go test -v ./tests/unit/user.service_test.go \
  		-coverpkg=./api/v1/services,./api/v1/repository \
  		-coverprofile=./tests/report/user/user.service.coverage.out
	go tool cover -html ./tests/report/user/user.service.coverage.out \
		-o ./tests/report/user/user.service.coverage.html

