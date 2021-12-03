#!/bin/bash
podman image exists localhost/remainders-backend:latest
if [[ $? == 1 ]]
then
  podman build -t localhost/remainders-backend:latest .
fi

podman pod exists remainders-backend
if [[ $? == 1 ]]
then
   podman pod create -p 8888:8888 --name remainders-backend
fi

podman container exists mongo-test
if [[ $? == 1 ]]
then
  podman run -d --name mongo-test --pod remainders-backend mongo
else
  podman restart mongo-test
fi
sleep 10

podman container exists backend-test
if [[ $? == 1 ]]
then
  podman run -d --name backend-test --pod remainders-backend -v "$PWD"/ssl:/var/app/ssl:Z -e PORT=8888 -e DB=remainders -e URI=localhost:27017 -e SSL_PUBLIC=./ssl/localhost.crt -e SSL_PRIVATE=./ssl/localhost.key -e ALLOWED_ORIGINS=http://localhost:3000 -e JWT_SECRET="topsecret" localhost/remainders-backend:latest
else
  podman restart backend-test
fi
