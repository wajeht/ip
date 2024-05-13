push:
	go test
	go fmt
	git add -A
	aicommits --type conventional
	git push --no-verify
