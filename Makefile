migrate-up:
	migrate -path=./migrations -database=$(MONDAY_HR_DB_DSN) up

migrate-down-1:
	migrate -path=./migrations -database=$(MONDAY_HR_DB_DSN) down 1

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=create_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ./migrations -seq $(name)