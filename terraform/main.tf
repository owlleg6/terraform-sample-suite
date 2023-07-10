module "resource_group" {
  source  = "data-platform-hq/resource-group/azurerm"
  version = "1.3.0"

  project  = var.project
  env      = var.env
  location = var.location

  custom_resource_group_name = var.custom_resource_group_name
  suffix                     = var.suffix

  tags = var.tags
}

module "log_analytics_ws" {
  count   = var.log_analytics_ws_enabled ? 1 : 0
  source  = "data-platform-hq/log-analytics-ws/azurerm"
  version = "1.2.0"

  project  = var.project
  env      = var.env
  location = var.location

  resource_group        = module.resource_group.name
  custom_workspace_name = var.custom_workspace_name
  retention_in_days     = var.analytics_retention_in_days

  tags = var.tags
}
