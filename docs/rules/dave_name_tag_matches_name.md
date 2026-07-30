# dave_name_tag_matches_name

Flags resources whose `Name` tag does not match their `name` attribute.

**Fixable:** Partially (see [Autofix](#autofix))

## Why

The `Name` tag is what shows up in the AWS console and in many third-party tools, while `name` is what Terraform actually sets. When the two drift apart, the console displays something different from the resource's real name, which is confusing during incident response and makes searching by name unreliable.

The `name` attribute is treated as the source of truth.

## When it applies

- Only resources with both a `name` attribute and a `tags` map are checked. Resources with no `name` argument (for example `aws_instance`, which relies entirely on the `Name` tag) are not checked — the rule has nothing to compare against.
- Resources with no `tags`, or tags with no `Name` key, are not checked.

## Autofix

When `tflint --fix` is run and the `Name` tag is a literal entry — directly in the tags map, or nested inside `merge(...)` — its value is replaced with the exact source text of the `name` attribute. This works whether `name` is a string literal or a reference (`var.name`, `local.name`, ...), since the fix mirrors source text rather than a resolved value.

```
$ tflint --fix
1 issue(s) found:

Warning: [Fixed] Name tag value "wrong" does not match the resource's name attribute value "web"; they must be identical (dave_name_tag_matches_name)
  on main.tf line 4:
   4:     Name = "wrong"
```

Autofix does not apply when:
- The `Name` entry can't be located precisely — for example `tags = var.tags`, where only the resolved value (not its source location) is known. The mismatch is still flagged, just not fixed.
- Both `name` and the `Name` tag are unresolvable references with different source text (see below) — a fix would require guessing which one is correct.

## Reference values

If `name` and the `Name` tag are both references that can't be resolved statically (e.g. `var.a`), the rule compares their source text instead of their values:

- Identical source (`name = var.a`, `Name = var.a`) is treated as consistent.
- Different source (`name = var.a`, `Name = var.b`) is flagged, since it usually indicates drift — though it's possible for two different expressions to resolve to the same value, in which case this is a false positive.
- If only one side is resolvable, the resource is skipped — comparing a known value against an unresolvable reference is too ambiguous to flag safely.

## Known limitations

- A `Name` key nested inside a non-static `merge()` argument (for example `merge(var.tags, ...)`) is not detected unless the whole `tags` expression still evaluates statically.
- `.tf.json` configuration is only checked via whole-map evaluation (no autofix), since there's no literal object syntax to locate a precise range in.

## Examples

```hcl
# ❌ Invalid (auto-fixed: Name becomes "web")
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    Name = "wrong"
  }
}

# ✅ Valid
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    Name = "web"
  }
}

# ✅ Valid — no name attribute, rule does not apply
resource "aws_instance" "this" {
  tags = {
    Name = "web"
  }
}
```
