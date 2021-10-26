.PHONY: all

run-e2e: docker-compose-cleanup
	docker-compose -f docker-compose.yml up -d
	./scripts/run-e2e.sh

docker-compose-cleanup:
	docker-compose -f docker-compose.yml down