migrateup:
	migrate -path api-server/db/migrations -database $(BLOCKCYPHER_DB_URL) up

migratedown:
	migrate -path api-server/db/migrations -database $(BLOCKCYPHER_DB_URL) down

sqlc:
	sqlc generate

.PHONY: migrateup sqlc migratedown