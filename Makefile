push:
	go test
	go fmt
	git add -A
	aicommits --type conventional
	git push --no-verify

test:
	go test

format:
	go fmt

update-db:
	curl -O https://git.io/GeoLite2-City.mmdb
