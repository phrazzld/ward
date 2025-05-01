# T014: Implement ci.Detector

## Task Description
Implement the ci.Detector interface which detects whether the code is running in a CI environment by checking environment variables.

## Approach

1. Review the existing interface in `internal/ci/detector.go`
2. Implement the `ci.Detector` in `internal/ci/ci.go` with:
   - A struct implementing the interface
   - The `IsCI()` method that checks for common CI environment variables:
     - Primary: `CI`
     - Specific: `GITHUB_ACTIONS`
     - Consider other common CI platforms as needed
3. Use a simple, efficient approach that works for the most common CI systems
4. Follow Go best practices from DEVELOPMENT_PHILOSOPHY.md
5. Error handling is minimal as this is a boolean function

## Implementation Details

1. Create a struct called `Detector` in the `ci` package
2. Implement the `IsCI()` method that:
   - Checks for the existence of the `CI` environment variable (set by most CI platforms)
   - Also checks for `GITHUB_ACTIONS` as specified in the task
   - Returns true if any of these variables are set to a non-empty value
3. Make the implementation thread-safe
4. Document the function with appropriate comments

## Testing Considerations (for T015)

1. Unit tests will need to:
   - Test with CI variables set (should return true)
   - Test with CI variables unset (should return false)
   - May require temporarily modifying environment variables during tests
