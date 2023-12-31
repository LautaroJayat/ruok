include .env
export $(shell sed 's/=.*//' .env)

migrate:
	go run cmd/main.go setupdb

seed:
	go run cmd/seed/main.go

build:
	go build -o ./ruok cmd/main.go

start-db:
	docker compose -f Dockercompose.dev.yml up --build

stop-db:
	docker compose -f Dockercompose.dev.yml down

start-db-tls:
	docker compose -f Dockercompose.dev.tls.yml up --build

stop-db-tls:
	docker compose -f Dockercompose.dev.tls.yml down

test:
	go test -p 1 -count=1 ./pkg/... -v

test-e2e:
	make build
	go test -p 1 -count=1 ./e2e/... -v
	make clean

run:
	go run cmd/main.go start

run-bin:
	./ruok start

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
	echo "hostssl all all 0.0.0.0/0 trust clientcert=verify-ca" > ssl/pg_hba.conf
	echo "hostnossl all all 0.0.0.0/0 reject" >> ssl/pg_hba.conf
	echo "local postgres user trust" >> ssl/pg_hba.conf
	sudo chown 70:70 ssl/server-key.pem

start-front:
	cd ruok-ui && npm run dev

build-front:
	rm -drf pkg/api/static
	mkdir pkg/api/static
	cd ruok-ui && npm run build
	 

clean:
	rm -f ./main
	rm -f ssl/*.pem
	rm -f *.pem
	rm -f ssl/*.conf
	rm -f e2e/*.log
	rm -f ./ruok