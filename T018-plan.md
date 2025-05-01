# T018: Implement git.Client adapter using util exec wrapper

## Task Description
Implement the git.Client interface in the `internal/git` package using the `ExecCommand` utility function to interact with the Git CLI. Additionally, implement placeholder file I/O in `.git/ward_placeholder` for tracking review results between pre-commit and post-commit hooks.

## Approach

1. Create a `GitClient` struct that implements the `git.Client` interface
2. Implement each method using the `util.ExecCommand` wrapper to call Git CLI commands
3. Use proper error handling with `errors.WrapExternalCmdError`
4. Implement placeholder file operations for the three placeholder-related methods
5. Add proper documentation for all methods
6. Ensure proper context propagation and cancellation handling

## Implementation Details

### GitClient Struct
- Create a struct that doesn't need any initialization fields
- Provide a constructor function `NewClient() *GitClient`

### Git CLI Command Implementations
1. **GetStagedDiff**: `git diff --cached`
2. **GetCommitMessage**: Read from `COMMIT_EDITMSG` file in git directory
3. **GetCurrentBranch**: `git symbolic-ref --short HEAD`
4. **GetHeadCommitHash**: `git rev-parse HEAD`
5. **GetGitDir**: `git rev-parse --git-dir`

### Placeholder File Operations
- For the placeholder file, use `.git/ward_placeholder`
- **WritePlaceholderHash**: Write the correlation ID to the placeholder file
- **ReadPlaceholderHash**: Read the correlation ID from the placeholder file
- **CleanPlaceholderHash**: Delete the placeholder file

### Error Handling
- All command execution errors should be wrapped with `errors.WrapExternalCmdError` (already done by `util.ExecCommand`)
- Add appropriate error handling for file operations

## Testing Considerations (for later T019)
- Will need to create temporary Git repositories for testing
- Should test all methods under normal and edge-case conditions
- Will need to verify placeholder file operations work correctly
