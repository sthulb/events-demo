build-GitHubConsumer:
	GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/main github-consumer/main.go
	cp github-consumer/bootstrap $(ARTIFACTS_DIR)/bootstrap
	chmod a+x $(ARTIFACTS_DIR)/bootstrap

build-EventsPublisher:
	GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/main events-publisher/main.go
	cp events-publisher/bootstrap $(ARTIFACTS_DIR)/bootstrap
	chmod a+x $(ARTIFACTS_DIR)/bootstrap

build-Issues:
	GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/main issues/main.go
	cp issues/bootstrap $(ARTIFACTS_DIR)/bootstrap
	chmod a+x $(ARTIFACTS_DIR)/bootstrap