output "resource_group" {
  value = {
    name     = module.resource_group.name
    id       = module.resource_group.id
    location = module.resource_group.location
  }
}

output "log_analytics_ws" {
  value = {
    name = try(module.log_analytics_ws[0].name, null)
    id   = try(module.log_analytics_ws[0].id, null)
  }
}
