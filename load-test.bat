@echo off
docker-compose run --rm k6 run /scripts/local-testing.js
@echo on