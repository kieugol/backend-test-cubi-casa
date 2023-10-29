## backend-dev-test
CubiCasa backend coding assignment

## Implement API with golang and postgres db

### Requirements

**Hub**, **Team**, and **Users** are in the order of hierarchy. A short description is as below:

- A Hub is an entity that associates Team depending on their geological location.
- A Team is an entity that associates Users based on their types.
- A User is an entity that holds the information of the human users.

**Todo**
- [ ] Implement a Create for each _hub, team, and user_.
- [ ] Implement a Search which will return _team and hub information_.
- [ ] Implement a Join for user into team, and team into hub (for simplicity: one user belongs to one team, one team belongs to one hub).
- [ ] Write the test cases
- [ ] Provide a SQL script which creates tables needed for the API.
- [ ] Good to use `docker/docker-compose` for local development setup(not mandatory)
- [ ] Good to provide the solution with security concern

### To submit this assignment
- Publish your code under your own Github account
- Send us an email that you've completed the task at openpositions@cubicasa.com with subject prefix **Backend Dev** and give a link to your repo
