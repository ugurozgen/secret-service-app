.SILENT:
vault:
	docker compose up -d

vault-down:
	docker compose down

run-app-locally:
	VAULT_ADDR=${VAULT_ADDR} \
	VAULT_TOKEN=${VAULT_TOKEN} \
	go run main.go 

run-app-container:
	docker run \
		-e VAULT_ADDR=${VAULT_ADDR} \
		-e VAULT_TOKEN=${VAULT_TOKEN} \
		-p 8080:8080 secret-service