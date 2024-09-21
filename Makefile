CONFIG_DIR=~/.code_deployer
.PHONY: install

install:
	@echo "Installing..."
	@mkdir -p $(CONFIG_DIR) 
	@touch $(CONFIG_DIR)/config.yaml
	@touch $(CONFIG_DIR)/service_account.yaml
	@echo "project_id: \"\"" > $(CONFIG_DIR)/config.yaml
	@echo "subscription_id: \"\"" >> $(CONFIG_DIR)/config.yaml
	@echo "docker_repo: \"\"" >> $(CONFIG_DIR)/config.yaml
	@echo "Please fill in the config file at $(CONFIG_DIR)/config.yaml"
	@echo "Add service account to $(CONFIG_DIR)/service_account.json"
	