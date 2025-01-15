# Makefile for initializing nodes and validators

# Variables
SCRIPTS_DIR=scripts

.PHONY: init init-nodes assign-validators create-accounts initialize-env generate-transactions

init: init-nodes assign-validators create-accounts initialize-env generate-transactions
	@echo "Initialization complete."

init-nodes:
	@echo "Initializing nodes..."
	bash $(SCRIPTS_DIR)/1_init_nodes.sh

assign-validators:
	@echo "Assigning validators..."
	bash $(SCRIPTS_DIR)/2_assign_validator.sh

create-accounts:
	@echo "Creating accounts..."
	bash $(SCRIPTS_DIR)/3_create_account.sh

initialize-env:
	@echo "Initializing environment and persistent peers..."
	bash $(SCRIPTS_DIR)/91_init.sh

generate-transactions:
	@echo "Generating transactions..."
	bash $(SCRIPTS_DIR)/4_generate_transactions.sh
