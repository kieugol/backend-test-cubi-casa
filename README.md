## backend task of CubiCasa
- [x] Implement a Create for each _hub, team, and user_.
- [x] Implement a Search which will return _team and hub information_.
- [x] Implement a Join for user into team, and team into hub (for simplicity: one user belongs to one team, one team belongs to one hub).
- [x] Write the test cases
- [x] Provide a SQL script which creates tables needed for the API.
- [x] Good to use `docker/docker-compose` for local development setup(not mandatory)
- [ ] Good to provide the solution with security concern

## Setup and test
1. **Requirement**
    * install `Docker` latest or stable version
    * Install `Make` package latest or stable version
    * Install `Postman`

2. **Install API and PostgreSQL DB on docker and setup ENV**
    1. CD to root folder `backend-test-cubi-casa`
    2. Open file `config/development.yaml` find config `database.host`
       change to your ip local
    3. Type `make install`
    4. Type `make start`  
       _(Please check carefully output same as below to make sure db and api stared successfully)_  
       `Creating backend-test-cubi-casa_db_1 ... done`  
       `Creating rest-api                    ... done`
    5. Type `make logs`  
       Output success display as sample output
       `... stdout: [GIN-debug] Listening and serving HTTP on :9090`

3. **Migrate DB** (Make sure step.2 have been completed)
   1. CD to root folder `backend-test-cubi-casa`
   2. Type `make ssh-app` (go in to container)
   3. Type `go run migrate/migrate.go`  
      see output complete ` 👍 Migration complete `


4. **Execute test with TODO task**
    - Prepare API: Using postman
        1. URI: `http://localhost:9090/`
        2. Header request: _data any(just make it as sample)_  
           `X-PLATFORM:IOS`  
           `X-DEVICE-TYPE:phone`  
           `X-DEVICE-ID:620005a5-1305-4668-9fb2-3ba250a57ab9`  
           `X-LANG:en`  
           `X-CHANNEL:2`
    - TODO task
      - Create for each hub, team, and user:
        - Api create hub: `v1/hubs`
          - **Method**: POST
          - **params**:  ```{
            "name": "hub9",
            "location": "HCM city",
            "description": "HCM city"
            } ```
        - Api create team: `v1/teams`
          - **Method**: POST
          - **params**:  ```{
            "name": "team-go",
            "hub_id": 1,
            "description": "team go"
            } ```
        - Api create user: `v1/users`
            - **Method**: POST
            - **params**:  ```{
              "name": "krol",
              "email": "krol@gmail.com",
              "phone": "+84917474772",
              "address": "HCM, City, p13, qTB",
              "team_id": 2
              } ```
      - Implement a Search which will return team and hub information  
          - Api search hub: `v1/teams?id=%d&name=%s`
              - **Method**: GET
              - **params**:  ```{
                "id": 1,
                "name": "krol"
                } ```
      - Implement a Join for user into team, and team into hub
        - Api search hub: `v1/users?name=%s&email=%s`
            - **Method**: GET
            - **params**:  ```{
              "name": "krol",
              "email": "kieu@"
              } ```    


5. **Execute Unit test**  (Using `make` command )
     1. CD to root folder `backend-test-cubi-casa`.  
     2. Type `make ssh-app` 
        - **Test Services** 
           1. Unit Test `hub.service.go` and `hub.repository.go`
               - Type `make test-hub-service` 
               - See report html at folder  `tests/report/hub/hub.service.coverage.html`
           2. Unit Test`team.service.go` and `team.repository.go`
              - Type `make test-team-service`
              - See report html at folder  `tests/report/team/team.service.coverage.html`
           3. Unit Test`user.service.go` and `user.repository.go`
              - Type `make test-user-service`
              - See report html at folder  `tests/report/user/user.service.coverage.html`
        - **Test Controller**
            1. Unit Test`user.controller.go`
               - Type `make test-user-controller`
               - See report html at folder  `tests/report/user/user.controller.coverage.html`
            2. Unit Test`team.controller.go` // TODO  
            3. Unit Test`hub.controller.go`  // TODO
