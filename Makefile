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
	rm -f GeoLite2-City.mmdb
	wget https://git.io/GeoLite2-City.mmdb

fix-git:
	git rm -r --cached .
	git add .
	git commit -m "chore: untrack files in .gitignore"
