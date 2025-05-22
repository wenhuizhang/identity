.PHONY: do_generate_proto do_stop_db do_start_db do_stop_network do_start_network do_start_node do_stop_node

do_generate_proto:
	cd scripts && ./buf-generate.sh
	@echo "Generated proto files"

do_generate_node_sdk:
	chmod +x scripts/node/generate.sh
	./scripts/node/generate.sh
	@echo "Generated Node SDK"

do_stop_db:
	@./deployments/scripts/identity/stop_db.sh
	@echo "CouchDB stopped"

do_start_db:
	@./deployments/scripts/identity/launch_db.sh
	@echo "CouchDB started at PORT 5984"

do_stop_network:
	@./deployments/scripts/network/stop_network.sh
	@echo "Identity network stopped"

do_start_network:
	@./deployments/scripts/network/launch_network.sh
	@echo "Identity network started"

do_start_node:
	@./deployments/scripts/identity/launch_node.sh
	@echo "Node started at http://localhost:4000"

do_stop_node:
	@./deployments/scripts/identity/stop_node.sh
	@echo "Node stopped"

generate_proto: do_generate_proto

generate_node_sdk: do_generate_node_sdk

stop_node: do_stop_node do_stop_db do_stop_network
start_node: stop_node do_start_network do_start_db do_start_node
