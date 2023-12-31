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
test-hub-service:
	go test -v ./tests/unit/hub.service_test.go \
  		-coverpkg=./api/v1/services,./api/v1/repository \
  		-coverprofile=./tests/report/hub/hub.service.coverage.out
	go tool cover -html ./tests/report/hub/hub.service.coverage.out \
		-o ./tests/report/hub/hub.service.coverage.html
test-team-service:
	go test -v ./tests/unit/team.service_test.go \
  		-coverpkg=./api/v1/services,./api/v1/repository \
  		-coverprofile=./tests/report/team/team.service.coverage.out
	go tool cover -html ./tests/report/team/team.service.coverage.out \
		-o ./tests/report/team/team.service.service.coverage.html
test-user-service:
	go test -v ./tests/unit/user.service_test.go \
  		-coverpkg=./api/v1/services,./api/v1/repository \
  		-coverprofile=./tests/report/user/user.service.coverage.out
	go tool cover -html ./tests/report/user/user.service.coverage.out \
		-o ./tests/report/user/user.service.coverage.html
test-user-controller:
	go test -v ./tests/unit/user.controller_test.go \
  		-coverpkg=./api/v1/controllers \
  		-coverprofile=./tests/report/user/user.controller.coverage.out
	go tool cover -html ./tests/report/user/user.controller.coverage.out \
		-o ./tests/report/user/user.controller.coverage.html

