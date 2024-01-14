
### Description
Simple project for a basic crowdfunding initiative sight.  With a backend written up in `Go` using the `Gin` framework.  Frontend setup based on `React` using `NextJS` framework.  Done up with basic `env` files and values to work out of the box for demo purposes.  Bundled together with telemetry feeding into a `Grafana` instance setup
  
Postman collection and environment provided for testing in the `postman` directory

### Running Locally
 - [ ] Clone repo locally
 - [ ] Run `docker compose up --build -d`
 - [ ] Seed data
	 - [ ] Go to seed directory with `cd api/seed`
	 - [ ] Run `go run .`

### Stopping
To tear down environment run `docker compose down`

### Backend
Generic REST API backend written up in `Go` with middleware authentication connecting up to a `MySQL` database backend leveraging `Gorm` framework.

`postman` directory contains some files for available requests and environment setup.  Once imported into Postman, these will setup the environment and add a collection with the endpoints.  All endpoints have tests, and the collection can be run with tests.  GUID for last entered entities will be saved for use with corresponding delete or patch routes.

### Frontend
No real progress beyond boilerplate setup

### Grafana
Default login the first time will be `admin` for username and password.  Should have data sources provisioned out of the box to wire up to view traces. On the side nav, select `Explore`, make sure that `Tempo` is selected in the dropdown box. Set the query type to `Search` and traces should be viewable providing requests have been sent to the backend.

### Viewing
- [Frontend](http://localhost:3000/)
- [Grafana Dashboard](http://localhost:3004/)