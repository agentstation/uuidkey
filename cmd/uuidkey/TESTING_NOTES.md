# CLI Testing Notes

## Current Test Coverage
- Library: 96.2%
- CLI: 80.1%

## Known Issues

### Cobra State Pollution (FIXED)
Tests were failing when run together due to Cobra retaining command arguments between test executions. This was resolved by creating fresh command instances for each test run using factory functions in `commands.go`.

The solution:
1. Created `newRootCommand()` and similar factory functions for each command
2. Modified `executeCommand()` test helper to use fresh instances
3. Added `setCurrentCommand()` to ensure output function uses the test command instance

### Other Test Issues (FIXED)
1. ✓ Key validation updated to accept 28 character keys
2. ✓ Output function updated to use current command instance
3. ✓ Test expectations updated to match current output formats
4. ✓ Crockford validation updated to match spec (Q is valid, I/L/O/U excluded)

## Lessons Learned
1. Cobra is designed for single-execution CLI apps, not repeated test runs
2. Global command instances should be avoided in testable code
3. Factory functions provide clean isolation for command testing
4. The library may generate characters outside strict Crockford spec (e.g., Q)