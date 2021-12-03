#!/bin/sh
k3d cluster create --port 8082:20080@agent:0 -p 8081:80@loadbalancer --agents 2
docker exec k3d-k3s-default-agent-0 mkdir /tmp/kube