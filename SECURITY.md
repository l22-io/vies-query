# Security Policy

## Supported Versions

The following versions of VIES Query receive security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | Yes |
| < 1.0   | No  |

We recommend always using the latest stable release for the best security posture.

## Security Features

VIES Query implements several security measures:

### Network Security
- **TLS 1.2+**: All HTTPS connections enforce TLS 1.2 or higher
- **Certificate Validation**: Proper certificate chain validation for VIES API
- **Timeout Protection**: Configurable timeouts prevent hanging connections
- **User Agent**: Proper identification in HTTP requests

### Input Validation
- **VAT Number Sanitization**: Input validation for VAT number formats
- **Country Code Validation**: Strict validation of EU country codes
- **Format Validation**: Regex-based validation before API calls

### Error Handling
- **No Information Disclosure**: Error messages don't expose sensitive data
- **Structured Errors**: Typed error handling prevents information leakage
- **Safe Logging**: Verbose mode logs are sanitized

### Build Security
- **Reproducible Builds**: Consistent build artifacts
- **Dependency Scanning**: Automated vulnerability scanning in CI
- **Static Analysis**: Security-focused static analysis tools
- **Supply Chain**: Minimal dependencies, Go standard library preferred

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in VIES Query, please report it responsibly.

### How to Report

1. **DO NOT** create a public GitHub issue for security vulnerabilities
2. Email us at: **security@l22.io**
3. Include the following information:
   - Description of the vulnerability
   - Steps to reproduce the issue
   - Potential impact assessment
   - Suggested fix (if available)

### What to Include

Please provide as much detail as possible:

- **Component**: Which part of the application is affected
- **Version**: What version(s) are affected
- **Environment**: Operating system, architecture, build method
- **Reproduction**: Detailed steps to reproduce the vulnerability
- **Impact**: Your assessment of the security impact
- **Fix**: Any potential fixes or mitigations you've identified

### Response Timeline

We will respond to security reports according to the following timeline:

- **24 hours**: Acknowledgment of receipt
- **72 hours**: Initial assessment and triage
- **7 days**: Detailed response with our analysis
- **30 days**: Security fix released (for confirmed vulnerabilities)

### Disclosure Policy

We follow responsible disclosure practices:

1. **Investigation**: We'll investigate and validate the reported vulnerability
2. **Fix Development**: If confirmed, we'll develop and test a fix
3. **Coordinated Release**: We'll coordinate the release with the reporter
4. **Public Disclosure**: After the fix is released, we'll publicly disclose the vulnerability
5. **Credit**: We'll credit the reporter (unless they prefer to remain anonymous)

## Security Best Practices

When using VIES Query in production:

### Network Security
- **Firewall Rules**: Configure firewall rules to allow only necessary traffic
- **Proxy Configuration**: If using a proxy, ensure it supports HTTPS
- **Network Monitoring**: Monitor network traffic for anomalies

### Operational Security
- **Regular Updates**: Keep VIES Query updated to the latest version
- **Access Control**: Limit access to the binary and configuration
- **Log Monitoring**: Monitor logs for suspicious activity
- **Rate Limiting**: Respect VIES service rate limits

### CI/CD Security
- **Secure Pipelines**: Use secure CI/CD practices
- **Binary Verification**: Verify checksums of downloaded binaries
- **Dependency Scanning**: Scan for vulnerabilities in build environments
- **Secrets Management**: Properly manage any secrets or API keys

### Development Security
- **Code Review**: All code changes should be reviewed
- **Security Testing**: Include security testing in your test suite
- **Input Validation**: Validate all inputs in your applications
- **Error Handling**: Handle errors securely without information disclosure

## Known Security Considerations

### VIES Service Dependencies
- **External Service**: VIES Query depends on the external European Commission VIES service
- **Network Reliability**: Service availability depends on network connectivity
- **Data Privacy**: Company information is retrieved from a third-party service
- **Rate Limiting**: The VIES service has rate limits that could impact availability

### Data Handling
- **No Data Storage**: VIES Query doesn't store VAT numbers or company information
- **Logging**: Verbose mode may log request/response data locally
- **Memory**: Company information is temporarily held in memory during processing
- **Network Traffic**: VAT numbers are transmitted to the VIES service over HTTPS

## Security Testing

We use several security testing approaches:

### Automated Testing
- **Static Analysis**: Gosec security scanner in CI
- **Dependency Scanning**: Nancy vulnerability scanner
- **Code Quality**: golangci-lint with security-focused rules
- **Fuzzing**: Property-based testing for input validation

### Manual Testing
- **Code Review**: Security-focused code reviews
- **Penetration Testing**: Regular security assessments
- **Network Analysis**: Network traffic analysis and monitoring

## Vulnerability Management

### Internal Process
1. **Detection**: Through automated scanning and manual review
2. **Assessment**: Evaluate severity and impact
3. **Prioritization**: Based on CVSS scores and exploitability
4. **Remediation**: Develop and test fixes
5. **Deployment**: Release security updates
6. **Communication**: Notify users of security updates

### Third-Party Dependencies
- **Monitoring**: Continuous monitoring of Go standard library security advisories
- **Updates**: Prompt updates when security fixes are available
- **Minimal Dependencies**: We minimize external dependencies to reduce attack surface

## Security Contact

For security-related questions or concerns:

- **Email**: security@l22.io
- **Response Time**: Within 24 hours for urgent security matters
- **PGP Key**: Available upon request for encrypted communications

## Acknowledgments

We appreciate security researchers who help keep our users safe:

- Security researchers who have reported vulnerabilities will be listed here
- We follow responsible disclosure and provide credit where appropriate
- Bug bounty information (if applicable) will be posted here

---

**Note**: This security policy is subject to change. Please check back regularly for updates.
