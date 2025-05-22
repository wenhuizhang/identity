.PHONY: do_generate_proto do_generate_node_sdk do_start_node_dev do_stop_node_dev

do_generate_proto:
	cd scripts && ./buf-generate.sh
	@echo "Generated proto files"

do_generate_node_sdk:
	chmod +x scripts/node/generate.sh
	./scripts/node/generate.sh
	@echo "Generated Node SDK"

do_start_node_dev:
	@./deployments/scripts/identity/launch_node_dev.sh
	@echo "Postgres started at :5984"
	@echo "Node started at :4000"

do_stop_node_dev:
	@./deployments/scripts/identity/stop_node_dev.sh
	@echo "Node stopped"
	@echo "Postgres stopped"

generate_proto: do_generate_proto

generate_node_sdk: do_generate_node_sdk

stop_node_dev: do_stop_node_dev
start_node_dev: do_start_node_dev
