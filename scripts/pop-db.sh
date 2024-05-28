#!/bin/sh
docker cp $1 store-db:/docker-entrypoint-initdb.d/dump.sql
docker exec store-db psql -U user -d store -f docker-entrypoint-initdb.d/dump.sql