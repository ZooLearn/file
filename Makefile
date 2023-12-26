.PHONY: gen-ent
gen-ent:
	go run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./ent/template/*.tmpl" ./ent/schema --feature sql/execquery,intercept

.PHONY: migrate-up
migrate-up:

.PHONY: migrate-down
migrate-down:

.PHONY: test
test:
	go test ./...

.PHONY: gen-domain-test
gen-domain-test:
	mockery --dir=domain --output=domain/mocks --outpkg=mocks --all

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run

.PHONY: go-vet
go-vet:
	go vet ./...

.PHONY: pre-commit
pre-commit: golangci-lint go-vet test
