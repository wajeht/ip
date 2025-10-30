commit:
	@git add -A
	@git auto

push:
	@go test ./...
	@go fmt ./...
	@git add -A
	@git auto
	@git push --no-verify

test:
	@go test ./...

build:
	@go build -o ip .

run:
	@go run .

format:
	@go fmt ./...

update-db:
	@rm -f GeoLite2-City.mmdb
	@wget https://git.io/GeoLite2-City.mmdb

fix-git:
	@git rm -r --cached .
	@git add .
	@git commit -m "chore: untrack files in .gitignore"

docker:
	@docker build -f Dockerfile -t ip .

docker-run:
	@docker run --rm -p 80:80 ip
