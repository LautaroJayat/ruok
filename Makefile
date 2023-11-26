include .env
export $(shell sed 's/=.*//' .env)

migrate:
	go run cmd/migrate/main.go

seed:
	go run cmd/seed/main.go

build:
	go build cmd/scheduler/main.go

start-db:
	docker compose -f Dockercompose.dev.yml up --build

stop-db:
	docker compose -f Dockercompose.dev.yml down

test:
	go test -count=1 ./...

run:
	./main

gen-ssl-conf:
	rm -f ssl/*.pem
	rm -f ssl/*.conf
	openssl req -x509 -newkey rsa:4096 -keyout ssl/ca-key.pem -out ssl/ca-cert.pem -days 365 -nodes -subj "/CN=BackendLabs/O=BackendLabs"
	openssl genpkey -algorithm RSA -out ssl/server-key.pem -aes256 -pass pass:serverpass
	openssl req -new -key ssl/server-key.pem -out ssl/server-csr.pem -subj "/CN=BackendLabs" -passin pass:serverpass
	openssl x509 -req -in ssl/server-csr.pem -CA ssl/ca-cert.pem -CAkey ssl/ca-key.pem -out ssl/server-cert.pem -days 365
	openssl genpkey -algorithm RSA -out ssl/_client-key.pem -aes256 -pass pass:clientpass
	openssl req -new -key ssl/_client-key.pem -out ssl/client-csr.pem -subj "/CN=BackendLabs" -passin pass:clientpass
	openssl x509 -req -in ssl/client-csr.pem -CA ssl/ca-cert.pem -CAkey ssl/ca-key.pem -out ssl/client-cert.pem -days 365
	openssl pkey -in ssl/_client-key.pem -out ssl/client-key.pem -traditional -passin pass:clientpass
	echo "hostssl db1 all 0.0.0.0/0 cert clientcert=verify-full" > ssl/pg_hba.conf
	echo "local postgres user trust" >> ssl/pg_hba.conf
	sudo chown 70:70 ssl/server-key.pem

clean:
	rm -f ./main
	rm -f ssl/*.pem
	rm -f *.pem
	rm -f ssl/*.conf
