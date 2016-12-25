# route53-auto-ddns

- `route53-auto-ddns` gets current machine's external ip and update with provided Route53 record.

## Usage

```
./route53-auto-ddns [domain]
```

## Constraints

- AWS credentials are read from environment variables and then `default` profile(`~/.aws/credentials`)
- This may create a new record when the provided not found
- TTL is set to 3600