.PHONY: do_generate_proto do_stop_mongo do_start_mongo do_stop_network do_start_network do_start_node do_stop_node

do_generate_proto:
	cd scripts && ./buf-generate.sh
	@echo "Generated proto files"

do_stop_mongo:
	./deployments/scripts/identity/stop_mongo.sh
	@echo "Mongo stopped"

do_start_mongo:
	./deployments/scripts/identity/launch_mongo.sh
	@echo "Mongo started at PORT 27017"

do_stop_network:
	./deployments/scripts/network/stop_identity_network.sh
	@echo "Identity network stopped"

do_start_network:
	./deployments/scripts/network/launch_identity_network.sh
	@echo "Identity network started"

do_start_node:
	./deployments/scripts/identity/launch_node.sh
	@echo "Node started at http://localhost:4000"

do_stop_node:
	./deployments/scripts/identity/stop_node.sh
	@echo "Node stopped"

generate_proto: do_generate_proto

stop_mongo: do_stop_mongo
start_mongo: do_stop_mongo do_start_mongo

stop_network: do_stop_network
start_network: do_stop_network do_start_network

start_node: do_start_node
stop_node: do_stop_node
