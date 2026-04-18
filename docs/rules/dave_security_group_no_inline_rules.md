# dave_security_group_no_inline_rules

Flags inline `ingress` and `egress` blocks on `aws_security_group`.

## Why

Inline rules cause conflicts when multiple Terraform configurations or manual changes modify the same security group. Standalone `aws_vpc_security_group_ingress_rule` and `aws_vpc_security_group_egress_rule` resources are individually manageable and produce cleaner plans.

## Examples

```hcl
# ❌ Invalid — inline ingress
resource "aws_security_group" "this" {
  name = "app-sg"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ✅ Valid — standalone rule
resource "aws_security_group" "this" {
  name = "app-sg"
}

resource "aws_vpc_security_group_ingress_rule" "https" {
  security_group_id = aws_security_group.this.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
}
```
