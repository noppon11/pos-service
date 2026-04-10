output "app_name" {
  description = "Application name"
  value       = var.app_name
}

output "environment" {
  description = "Environment name"
  value       = var.environment
}

output "resource_id" {
  description = "Example null resource ID"
  value       = null_resource.example.id
}