#! /bin/bash

# Set environment variables for local development

export AZURE_TENANT_ID=$(az account show --query tenantId -o tsv)
export AZURE_SUBSCRIPTION_ID=$(az account show --query id -o tsv)