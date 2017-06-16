docker rm -f container-pool-apiserver
docker run -v /etc/localtime:/etc/localtime -d -p 8080:8080 -ti --name container-pool-apiserver   --link container-pool-db:db --link container-pool-proxy:proxy li41898/container-pool ./controller -D server --rethinkdb-addr db:28015 -d tcp://proxy:2375
