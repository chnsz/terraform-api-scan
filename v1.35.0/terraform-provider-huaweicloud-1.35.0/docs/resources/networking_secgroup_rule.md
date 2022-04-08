---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud_networking_secgroup_rule

Manages a Security Group Rule resource within HuaweiCloud. This is an alternative
to `huaweicloud_networking_secgroup_rule_v2`

## Example Usage

```hcl
resource "huaweicloud_networking_secgroup" "mysecgroup" {
  name        = "secgroup"
  description = "My security group"
}

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule" {
  security_group_id = huaweicloud_networking_secgroup.mysecgroup.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "334,466-468,8000"
  remote_ip_prefix  = "0.0.0.0/0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the security group rule resource. If
  omitted, the provider-level region will be used. Changing this creates a new security group rule.

* `security_group_id` - (Required, String, ForceNew) Specifies the security group id the rule should belong to. Changing
  this creates a new security group rule.

* `direction` - (Required, String, ForceNew) Specifies the direction of the rule, valid values are **ingress** or
  **egress**. Changing this creates a new security group rule.

* `ethertype` - (Required, String, ForceNew) Specifies the layer 3 protocol type, valid values are **IPv4** or **IPv6**.
  Changing this creates a new security group rule.

* `description` - (Optional, String, ForceNew) Specifies the supplementary information about the networking security
  group rule. This parameter can contain a maximum of 255 characters and cannot contain angle brackets (< or >).
  Changing this creates a new security group rule.

* `protocol` - (Optional, String, ForceNew) Specifies the layer 4 protocol type, valid values are **tcp**, **udp**,
  **icmp** and **icmpv6**. If omitted, the protocol means that all protocols are supported.
  This is required if you want to specify a port range. Changing this creates a new security group rule.

* `ports` - (Optional, String, ForceNew) Specifies the allowed port value range, which supports single port (80),
  continuous port (1-30) and discontinous port (22, 3389, 80) The valid port values is range form `1` to `65,535`.
  Changing this creates a new security group rule.

* `remote_ip_prefix` - (Optional, String, ForceNew) Specifies the remote CIDR, the value needs to be a valid CIDR (i.e.
  192.168.0.0/16). Changing this creates a new security group rule.

* `remote_group_id` - (Optional, String, ForceNew) Specifies the remote group id. Changing this creates a new security
  group rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 10 minute.

## Import

Security Group Rules can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_networking_secgroup_rule.secgroup_rule_1 aeb68ee3-6e9d-4256-955c-9584a6212745
```
