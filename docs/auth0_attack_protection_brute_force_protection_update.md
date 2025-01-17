---
layout: default
---
## auth0 attack-protection brute-force-protection update

Update brute force protection settings

### Synopsis

Update brute force protection settings.

```
auth0 attack-protection brute-force-protection update [flags]
```

### Examples

```
auth0 attack-protection brute-force-protection update
```

### Options

```
  -l, --allowlist strings   List of trusted IP addresses that will not have attack protection enforced against them. Comma-separated.
  -e, --enabled             Enable (or disable) brute force protection.
  -h, --help                help for update
  -a, --max-attempts int    Maximum number of unsuccessful attempts. (default 1)
  -m, --mode string         Account Lockout: Determines whether or not IP address is used when counting failed attempts. Possible values:
                            count_per_identifier_and_ip, count_per_identifier.
  -s, --shields strings     Action to take when a brute force protection threshold is violated. Possible values: block, user_notification. Comma-separated.
```

### Options inherited from parent commands

```
      --debug           Enable debug mode.
      --force           Skip confirmation.
      --format string   Command output format. Options: json.
      --no-color        Disable colors.
      --no-input        Disable interactivity.
      --tenant string   Specific tenant to use.
```

### SEE ALSO

* [auth0 attack-protection](auth0_attack_protection.md)	 - Manage attack protection settings
* [auth0 attack-protection brute-force-protection](auth0_attack_protection_brute_force_protection.md)	 - Manage brute force protection settings
