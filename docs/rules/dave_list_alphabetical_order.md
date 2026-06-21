# dave_list_alphabetical_order

Ensures the string elements of a list literal are sorted alphabetically, for a configurable set of attribute names.

**Fixable:** Yes (single-line lists only)

## Why

Alphabetical ordering makes a list easy to scan and produces predictable diffs — a new element slots into a known position rather than landing wherever the author happened to add it.

This rule is opt-in: there is no universally-correct set of lists to sort, and applying it everywhere would be noise. By naming an attribute in the configuration you are asserting that element order is not semantically significant for it — so the rule does nothing until you list the attributes you care about.

Lists that contain anything other than static strings (a variable reference, an interpolation, a function call) are left alone, because reordering could change behaviour or the intended order is ambiguous.

## Configuration

The rule is a no-op unless `attributes` is set.

```hcl
rule "dave_list_alphabetical_order" {
  enabled          = true
  attributes       = ["tags", "subnet_ids"]  # default: []
  case_insensitive = false                    # default: false
}
```

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `attributes` | `list(string)` | `[]` | Attribute names whose list values must be sorted. Matched anywhere, including inside nested blocks. |
| `case_insensitive` | `bool` | `false` | When `false`, elements are compared in byte order (capitals before lowercase). When `true`, comparison is on a lower-cased key. |

## Autofix

When run with `tflint --fix`, this rule reorders the elements of an unsorted list, preserving each element's original quoting and escaping.

Autofix only applies to single-line lists with no comments. Multiline lists and lists containing comments are flagged but not fixed, to avoid guessing indentation and to keep comments attached to the right element.

```
$ tflint --fix
1 issue(s) found:

Warning: [Fixed] List assigned to 'tags' is not sorted alphabetically. (dave_list_alphabetical_order)
  on main.tf line 2:
   2:   tags = ["b", "a", "c"]
```

For multiline lists, run `LC_ALL=C sort` to sort the elements alphabetically.

## Examples

Assuming `attributes = ["tags"]`:

```hcl
# ❌ Invalid — out of order (auto-fixable)
locals {
  tags = ["b", "a", "c"]
}

# ❌ Invalid — out of order, flagged but not auto-fixed (multiline)
locals {
  tags = [
    "b",
    "a",
  ]
}

# ✅ Valid
locals {
  tags = ["a", "b", "c"]
}

# ✅ Ignored — contains a non-static element
locals {
  tags = ["b", var.extra, "a"]
}
```
