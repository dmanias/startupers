# Makefile for setting up and managing the creative-narration project

# Variables
PAT ?= your_pat_here
GITHUB_USERNAME ?= your_github_username_here
HELM_VERSION ?= v3.0.0
SKAFFOLD_VERSION ?= latest

.PHONY: install-pre-commit setup-pre-commit install-helm install-skaffold login-ghcr skaffold-dev skaffold-run

# Install pre-commit hooks and related setup
install-pre-commit:
	@echo "Installing pre-commit..."
	@pip install pre-commit || sudo apt install python3-pre-commit

setup-pre-commit:
	@echo "Setting up pre-commit hooks..."
	@pre-commit install

# Install Helm
install-helm:
	@echo "Downloading and installing Helm $(HELM_VERSION)..."
	@curl -L https://get.helm.sh/helm-$(HELM_VERSION)-linux-amd64.tar.gz -o helm.tar.gz
	@tar -zxvf helm.tar.gz
	@mv linux-amd64/helm /usr/local/bin/helm
	@rm -rf linux-amd64 helm.tar.gz

# Install Skaffold
install-skaffold:
	@echo "Installing Skaffold..."
	@curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/$(SKAFFOLD_VERSION)/skaffold-linux-amd64 && chmod +x skaffold && sudo mv skaffold /usr/local/bin

# Login to GitHub Container Registry
login-ghcr:
	@echo "Logging into GitHub Container Registry..."
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USERNAME) --password-stdin

# Run Skaffold in development mode
skaffold-dev:
	@echo "Running Skaffold in development mode..."
	@skaffold dev

# Run Skaffold with --tail option
skaffold-run:
	@echo "Running Skaffold with --tail option..."
	@skaffold run --tail

# Default target to display help
help:
	@echo "Makefile for managing the creative-narration project"
	@echo ""
	@echo "Usage:"
	@echo "  make install-pre-commit   Install pre-commit hooks."
	@echo "  make setup-pre-commit     Setup pre-commit hooks."
	@echo "  make install-helm         Download and install Helm."
	@echo "  make install-skaffold     Install Skaffold."
	@echo "  make login-ghcr           Login to GitHub Container Registry."
	@echo "  make skaffold-dev         Run Skaffold in development mode."
	@echo "  make skaffold-run         Run Skaffold with --tail option."
	@echo "  make help                 Display this help message."
