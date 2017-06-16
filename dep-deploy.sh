docker rm -f container-pool-db
docker run -ti -d --restart=always --name container-pool-db -p 9090:8080 -p 28015:28015 rethinkdb
docker rm -f container-pool-proxy
docker run -ti -d -p 2375:2375 -v /var/run/docker.sock:/var/run/docker.sock --restart=always --name container-pool-proxy li41898/docker-proxy ./docker-proxy
