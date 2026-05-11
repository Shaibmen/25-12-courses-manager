#!/bin/sh

cd document-service

docker build -t metapohuism/document-service .

cd ..

docker compose up
