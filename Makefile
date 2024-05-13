push:
	go test
	go fmt
	git add -A
	./commit.sh
	git push --no-verify
