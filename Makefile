commit:
	npm run test
	npm run format
	npm run lint
	git add -A
	./commit.sh
	git push --no-verify
