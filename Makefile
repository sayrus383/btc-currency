start:
	docker build -t btc-currency-image .
	docker run -p 80:80 --name btc-currency-app btc-currency-image

test:
	go test ./... -v -short