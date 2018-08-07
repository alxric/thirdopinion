-include CONFIG
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git rev-parse --verify HEAD)
GOVER := $(shell go version | awk '{print $$3}' | cut -c 3-)

build:
	docker build --pull --rm \
	--build-arg RELEASE=$(VERSION) \
	--build-arg BRANCH=$(BRANCH) \
	--build-arg COMMIT=$(COMMIT) \
	--build-arg GOVER=$(GOVER) \
	--tag=$(APPLICATION):$(VERSION)-$(USER) \
	--label build.user=$(USER) \
	--label build.name=$(APPLICATION) .
