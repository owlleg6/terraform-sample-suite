variable "project" {
  type        = string
  description = "Project name"
}

variable "env" {
  type        = string
  description = "Environment name"
}

variable "location" {
  type        = string
  description = "Azure location"
}

variable "custom_resource_group_name" {
  type        = string
  description = "Custom name for Resource Group"
  default     = null
}

variable "suffix" {
  type        = string
  description = "Optional suffix for resource group"
  default     = ""
}

variable "tags" {
  type        = map(any)
  description = "A mapping of tags to assign to the resource"
  default = {
    env = "testing"
  }
}

# Log Analytics
variable "log_analytics_ws_enabled" {
  type    = bool
  default = false
}

variable "custom_workspace_name" {
  type        = string
  description = "Custom name for Log Analytics"
  default     = null
}

variable "analytics_retention_in_days" {
  type        = number
  description = "Custom name for Log Analytics"
  default     = 45
}
