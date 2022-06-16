BINARY=engine

engine:
	go build -o ${BINARY}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

migrate:
	go run migrations/main.go

docker:
	docker build -t simple-jwt-go .

run:
	docker-compose up --build -d

stop:
	docker-compose down
	