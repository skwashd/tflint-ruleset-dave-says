plugin "dave-says" {
  enabled = true
  # version = "0.2.0"
  # source  = "github.com/skwashd/tflint-ruleset-dave-says"
}

rule "dave_cloudwatch_log_retention" {
  enabled        = true
  retention_days = 30
}

rule "dave_list_alphabetical_order" {
  enabled          = true
  attributes       = ["tags", "subnet_ids"]
  case_insensitive = false
}
