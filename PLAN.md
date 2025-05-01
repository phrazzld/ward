```markdown
# Plan: Refactor Ward from Shell Scripts to Go

## Chosen Approach (One‑liner)

Implement Ward as a modular Go CLI application using internal interfaces to abstract external dependencies (Git CLI, LLM CLI, File System), enabling high testability, maintainability, and adherence to our development philosophy while preserving existing functionality.

---

## Architecture Blueprint

### **Modules / Packages**

| Package           | Responsibility                                                                                                | Key Dependencies (Interfaces)           | Notes                                                                 |
| :---------------- | :------------------------------------------------------------------------------------------------------------ | :-------------------------------------- | :-------------------------------------------------------------------- |
| `cmd/ward`        | Main entry point, CLI command structure (Cobra), flag parsing, dependency injection setup, subcommand dispatch. | `core.Checker`, `core.LogUpdater`, `core.Initializer` | Orchestrates application flow.                                        |
| `internal/config` | Load and validate configuration (env vars primarily: log path, temp dir, philosophy path, LLM settings).        | -                                       | Provides `Config` struct. Fails fast on invalid config.               |
| `internal/core`   | Core business logic for `check`, `log`, `init`. Defines primary interfaces (`GitClient`, `LLMClient`, `Logger`). | `git.Client`, `llm.Client`, `log.Writer`, `ci.Detector`, `shim.Generator` | Pure logic, no direct external interaction.                           |
| `internal/git`    | `git.Client` interface definition and implementation using `os/exec` to call the Git CLI.                       | -                                       | Abstracts diff, commit message, hash, branch operations.              |
| `internal/llm`    | `llm.Client` interface definition and implementation using `os/exec` to call the `claude` CLI.                  | -                                       | Abstracts LLM review request/response. Designed for future API swap. |
| `internal/log`    | `log.Writer` interface definition and implementation for structured JSON Lines logging to `.ward-warnings.log`. | -                                       | Handles log entry creation, finding, and updating by placeholder hash. |
| `internal/shim`   | `shim.Generator` interface and implementation to generate hook scripts for managers (pre-commit, husky).      | -                                       | Generates minimal wrappers calling `ward check`/`log`.                 |
| `internal/ci`     | `ci.Detector` interface and implementation to detect CI environments (e.g., `CI`, `GITHUB_ACTIONS`).          | -                                       | Allows skipping hooks in CI.                                          |
| `internal/errors` | Defines custom error types/variables (e.g., `ErrCheckFailed`, `ErrLogUpdateNotFound`).                        | -                                       | Standardizes error handling.                                          |
| `internal/types`  | Shared domain types (`ReviewRequest`, `ReviewResult`, `LogEntry`, `Config`).                                  | -                                       | Centralizes core data structures.                                     |
| `internal/util`   | Utility functions (correlation ID generation, file system helpers, safe command execution).                   | -                                       | Common helper functions.                                              |

### **Public Interfaces / Contracts**

```go
// internal/git/client.go
type Client interface {
    GetStagedDiff(ctx context.Context) (string, error)
    GetCommitMessage(ctx context.Context, commitMsgPath string) (string, error)
    GetCurrentBranch(ctx context.Context) (string, error)
    GetHeadCommitHash(ctx context.Context) (string, error)
    GetGitDir(ctx context.Context) (string, error)
    WritePlaceholderHash(ctx context.Context, placeholder string) error
    ReadPlaceholderHash(ctx context.Context) (string, error)
    CleanPlaceholderHash(ctx context.Context) error
}

// internal/llm/client.go
type Client interface {
    AnalyzeChanges(ctx context.Context, req types.ReviewRequest) (*types.ReviewResult, error)
}

// internal/log/writer.go
type Writer interface {
    WriteEntry(ctx context.Context, entry types.LogEntry) error
    UpdateEntryCommitHash(ctx context.Context, placeholderHash, actualHash string) error
}

// internal/shim/generator.go
type Generator interface {
    Generate(ctx context.Context, hookManager string, hookType string) (string, error) // Returns script content
}

// internal/ci/detector.go
type Detector interface {
    IsCI() bool
}

// internal/core interfaces (used by cmd)
type Checker interface {
    Check(ctx context.Context, commitMsgPath string) error // Returns ErrCheckFailed on LLM FAIL/WARN
}

type LogUpdater interface {
    UpdateLog(ctx context.Context) error
}

type Initializer interface {
    Initialize(ctx context.Context, hookManager string) error
}

// Shared types in internal/types
type ReviewRequest struct { /* Diff, CommitMsg, Philosophy */ }
type ReviewResult struct { /* Status (PASS/WARN/FAIL), Feedback, RawOutput */ }
type LogEntry struct { /* Timestamp, Status, CommitHash (placeholder/actual), Branch, CommitMsgHead, CorrelationID, Feedback */ }
type Config struct { /* LogFilePath, TempDir, PhilosophyPath, LLMProviderConfig */ }
```

### **Data Flow Diagram** (Mermaid)

```mermaid
graph TD
    subgraph User/System Interaction
        HookMgrPre[Hook Manager (pre-commit)] --> CLI_Check[ward check <commit_msg_file>]
        HookMgrPost[Hook Manager (post-commit)] --> CLI_Log[ward log]
        UserCLIInit[User CLI] --> CLI_Init[ward init --manager <name>]
    end

    subgraph Ward Go Application
        CLI_Check --> CmdCheck{cmd/check}
        CLI_Log --> CmdLog{cmd/log}
        CLI_Init --> CmdInit{cmd/init}

        CmdCheck --> CIDetect1[ci.Detector]
        CIDetect1 -- Not CI --> ConfigLoad1[config.Loader]
        ConfigLoad1 --> CoreCheck[core.Checker]
        CoreCheck -- Uses --> GitClient1[git.Client]
        CoreCheck -- Uses --> LLMClient1[llm.Client]
        CoreCheck -- Uses --> LogWriter1[log.Writer]
        CoreCheck -- Uses --> UtilCorrID[util.CorrelationID]

        CmdLog --> CIDetect2[ci.Detector]
        CIDetect2 -- Not CI --> ConfigLoad2[config.Loader]
        ConfigLoad2 --> CoreLogUpdate[core.LogUpdater]
        CoreLogUpdate -- Uses --> GitClient2[git.Client]
        CoreLogUpdate -- Uses --> LogWriter2[log.Writer]

        CmdInit --> ConfigLoad3[config.Loader]
        ConfigLoad3 --> CoreInit[core.Initializer]
        CoreInit -- Uses --> ShimGen[shim.Generator]

        CoreCheck -- Exit Code (0/1) --> CLI_Check
        CoreLogUpdate -- Exit Code (0/1) --> CLI_Log
        CoreInit -- Output/Exit Code --> CLI_Init
    end

    subgraph External Dependencies (Abstracted by Interfaces)
        GitClient1 --> GitCLI[(git CLI via os/exec)]
        GitClient2 --> GitCLI
        LLMClient1 --> ClaudeCLI[(claude CLI via os/exec)]
        LogWriter1 --> LogFile[FS: .ward-warnings.log]
        LogWriter2 --> LogFile
        GitClient1 --> PlaceholderFile[FS: .git/ward_placeholder]
        GitClient2 --> PlaceholderFile
        ShimGen --> HookScripts[FS: Hook Scripts]
    end
```

### **Error & Edge‑Case Strategy**

-   **Error Types:** Use `internal/errors` for specific, checkable errors (e.g., `errors.ErrCheckFailed`, `errors.ErrLogUpdateNotFound`, `errors.ErrExternalCmdFailed`, `errors.ErrConfigInvalid`). Standard Go errors for general issues.
-   **Error Propagation:** Wrap errors with context using `fmt.Errorf("module: operation: %w", err)` as they cross package boundaries.
-   **Exit Codes:**
    -   `ward check`: `0` (PASS), `1` (WARN/FAIL reported by LLM), `>1` (Internal error - config, Git, LLM CLI failure, etc.).
    -   `ward log`: `0` (Success), `1` (Internal error - failed to read placeholder, update log, etc.).
    -   `ward init`: `0` (Success), `1` (Error - invalid manager, write failure).
-   **External Failures:** `os/exec` calls must check exit codes and capture stderr. Wrap these into `errors.ErrExternalCmdFailed`.
-   **Empty Diff:** `ward check` exits `0` with an info message. No LLM call.
-   **CI Environment:** Hooks (`check`, `log`) exit `0` immediately if `ci.Detector.IsCI()` is true.
-   **File I/O:** Handle errors robustly (permissions, not found, disk full) for log file and placeholder file. Log updates must be atomic or idempotent if possible (JSON Lines append helps).
-   **LLM Parsing:** Defensively parse LLM output. If format is unexpected, treat as an internal error, log details, and potentially exit `>1`.
-   **Placeholder Handling:** `ward check` writes placeholder; `ward log` reads and cleans it. Handle cases where placeholder is missing during `log` (log warning, exit 0).

---

## Detailed Build Steps

1.  **Project Setup:** `go mod init`, directory structure, `golangci-lint` config, `Makefile` for common tasks (lint, test, build).
2.  **Types Definition:** Define structs in `internal/types` (`ReviewRequest`, `ReviewResult`, `LogEntry`, `Config`).
3.  **Error Definition:** Define sentinel errors/types in `internal/errors`.
4.  **Interfaces:** Define all interfaces (`git.Client`, `llm.Client`, `log.Writer`, `shim.Generator`, `ci.Detector`, `core.*`).
5.  **Config Implementation:** Implement `internal/config` loader (env vars first). Add unit tests.
6.  **CI Detector Implementation:** Implement `internal/ci`. Add unit tests.
7.  **Utility Implementation:** Implement `internal/util` (correlation ID, safe `os/exec` wrapper). Add unit tests.
8.  **Git Client Implementation:** Implement `internal/git` using the `os/exec` wrapper from `util`. Focus on robust command execution, output parsing, error handling, and placeholder file I/O. Add integration tests (requires test repo setup or careful `os/exec` mocking).
9.  **LLM Client Implementation:** Implement `internal/llm` using `os/exec` wrapper. Handle prompt construction, output parsing (PASS/WARN/FAIL), error handling. Add integration tests (mock `os/exec`).
10. **Log Writer Implementation:** Implement `internal/log` using `encoding/json` (JSON Lines). Handle file creation, append, read/search (for update), atomic update logic (e.g., read all, update in memory, write to temp, rename). Use file locking (`github.com/gofrs/flock` or similar). Add integration tests (real file system).
11. **Shim Generator Implementation:** Implement `internal/shim` for target hook managers (e.g., `pre-commit`). Add unit tests.
12. **Core Logic Implementation:** Implement `internal/core` (`check`, `log`, `init` logic). Inject dependencies via interfaces. Write extensive unit tests with mocks (e.g., `gomock`).
13. **CLI Commands (Cobra):** Implement `cmd/ward` subcommands. Parse flags, setup dependencies (DI container or manual wiring), call core logic methods, handle errors, print user feedback, set exit codes.
14. **E2E Testing:** Create basic E2E test script (bash/Go test) simulating git workflow (add, commit) triggering the compiled `ward` binary, verifying exit codes and log file content.
15. **Build & Release:** Configure GoReleaser for cross-platform static binary builds.
16. **Documentation:** Write Go doc comments. Update `README.md` (installation, usage, config, log format). Add architecture diagram.

---

## Testing Strategy

-   **Unit Tests:**
    -   **Scope:** `internal/core`, `internal/config`, `internal/ci`, `internal/util`, `internal/shim`, pure logic within adapters.
    -   **Mocks:** Use `gomock` (or similar) to mock interfaces (`git.Client`, `llm.Client`, `log.Writer`, etc.) injected into `internal/core`.
    -   **Goal:** Verify business logic correctness, edge cases (empty diff, errors from deps), configuration parsing. High coverage (>90%).
-   **Integration Tests:**
    -   **Scope:** Test adapter implementations (`internal/git`, `internal/llm`, `internal/log`) against their real (or simulated) external boundaries.
    -   **Mocks:** Mock `os/exec` calls for `git` and `llm` adapters to control external CLI behavior without running actual commands. Test `internal/log` against the real file system in a temporary directory.
    -   **Goal:** Verify correct interaction with external processes and the file system (command args, parsing output, file reads/writes/locking). Moderate coverage (>70%).
-   **E2E Tests:**
    -   **Scope:** Full CLI execution (`ward check`, `ward log`) triggered by simulated Git hooks in a test repository.
    -   **Mocks:** None (uses compiled binary, real Git CLI, mocked LLM CLI via script/path).
    -   **Goal:** Verify the complete workflow, including hook integration, CLI argument parsing, exit codes, and final log file state. Small number of critical path tests.
-   **What to Mock:** Strictly follow philosophy: mock **only** true external system boundaries or the interfaces abstracting them. In unit tests, mock the Go interfaces (`git.Client`, `llm.Client`, etc.). In integration tests, mock the `os/exec` calls *if* necessary, or use real file system.
-   **Coverage Targets:** Enforce >85% overall line coverage in CI. `internal/core` should aim for >95%.

---

## Logging & Observability

-   **Primary Log:** `.ward-warnings.log` (JSON Lines format) for structured review results (PASS/WARN/FAIL). Configurable via `WARD_LOG_FILE`.
    -   **Fields:** `timestamp` (RFC3339 UTC), `level` ("PASS", "WARN", "FAIL"), `commitHash` (placeholder -> actual), `branch`, `commitMsgHead` (first line), `correlationId`, `feedback` (LLM output), `provider` ("claude_cli").
-   **Operational Logs:** Log to `stderr` for operational info/debug/errors using `log/slog` (or `zerolog`).
    -   **Format:** Console-friendly during interactive runs, JSON if `WARD_LOG_FORMAT=json` env var is set.
    -   **Events:** App start/end, config loading, command execution, external calls (Git, LLM), file operations, errors encountered.
    -   **Fields:** `timestamp`, `level` (debug, info, warn, error), `msg`, `correlationId` (if applicable), `module`, `error` (if applicable), `duration_ms` (for external calls).
-   **Correlation ID:** Generate UUID at the start of `ward check`. Include in operational logs (stderr) and the review result entry (`.ward-warnings.log`). Propagate via `context.Context`.

---

## Security & Config

-   **Input Validation:**
    -   Sanitize/validate commit message file path argument.
    -   Limit size of diff/commit message passed to LLM (`llm.Client` responsibility).
    -   Validate configuration values (`internal/config`).
    -   Validate expected format of LLM response (`internal/llm`).
-   **Secrets Handling:**
    -   Currently none directly handled (relies on pre-configured `claude` CLI).
    -   **Future:** If LLM API keys are added, load **only** from environment variables (`WARD_LLM_API_KEY`) or secure secrets management. **Never** log secrets.
-   **Least Privilege:**
    -   Go binary needs read access to repo files/`.git`, write access to `.ward-warnings.log` and `.git/ward_placeholder`.
    -   Use `exec.CommandContext` without shell interpolation (`bash -c`) to prevent command injection when calling `git` and `claude`. Pass arguments directly.
-   **Dependency Security:** Use `govulncheck` in CI to scan for known vulnerabilities in dependencies.

---

## Documentation

-   **Code Self-Doc:** Standard Go doc comments (`//`) for all exported types, functions, constants, and interfaces. Explain *why* not just *what*.
-   **README.md:**
    -   Update Installation (Go install, binaries from releases).
    -   Update Usage (`ward check`, `ward log`, `ward init`).
    -   Update Hook Manager examples (`.pre-commit-config.yaml`, etc.) to call `ward`.
    -   Document Configuration (environment variables).
    -   Document `.ward-warnings.log` JSON Lines format.
    -   Include Architecture section with diagram link/embed.
-   **Diagrams:** Keep Mermaid diagrams (`ARCHITECTURE.md` or similar) updated.

---

## Risk Matrix

| Risk                                        | Severity | Mitigation                                                                                                                                     |
| :------------------------------------------ | :------- | :--------------------------------------------------------------------------------------------------------------------------------------------- |
| `claude` CLI dependency / output changes    | High     | Abstract via `llm.Client` interface. Robust parsing, defensive error handling. Integration tests mocking CLI output. Plan for future API swap. |
| Git command variations / edge cases         | Medium   | Use standard Git commands. Abstract via `git.Client`. Integration tests with diverse repo states (bare, submodules, etc.). Careful error parsing. |
| Log file corruption / race conditions       | Medium   | Use JSON Lines (append-friendly). Implement file locking (`flock`) for read/update operations in `log.Writer`. Atomic file replace for updates.  |
| Placeholder hash sync issues                | Medium   | Use dedicated file (`.git/ward_placeholder`). Include Correlation ID in file and log. Robust read/write/cleanup logic in `git.Client`.         |
| Performance worse than shell scripts        | Low      | Go startup is fast. `os/exec` has overhead but likely less than chained shell scripts. Profile if benchmarks show regression.                   |
| Security: Command Injection via `os/exec`   | Medium   | **Never** use shell interpolation (`bash -c`). Use `exec.CommandContext` with direct arguments. Sanitize any user input used in args.           |
| Overengineering / Scope Creep               | High     | Stick strictly to refactoring existing features. Defer new features (e.g., direct API, other LLMs) until post-refactor. Ruthless code review.    |
| Incomplete backward compatibility           | Medium   | Ensure CLI commands (`check`, `log`) accept compatible arguments/env vars where necessary. Provide clear migration docs.                      |
| Testing gaps (especially integration/E2E) | Medium   | Enforce coverage minimums. Prioritize integration tests for `git`, `llm`, `log` adapters. Implement core E2E workflow tests.                     |

---

## Open Questions

1.  Confirm the exact required behavior if `.ward-warnings.log` is missing during the `ward log` update (Current plan: Log warning, exit 0).
2.  Should PASS results be logged to `.ward-warnings.log`? (Current plan: No, matching script behavior, but easy to add via config flag later).
3.  Final list of hook managers `ward init` must support initially? (`pre-commit` assumed minimum).
4.  Is direct Claude API usage desired post-refactor, or is sticking with the CLI acceptable long-term? (Affects priority of swapping `llm.Client` implementation).
5.  Any specific constraints or requirements regarding the content/format of the Development Philosophy summary passed to the LLM? (Assume plain text path for now).
6.  Are there size limits or specific content types within diffs that should be explicitly excluded from the LLM review?
```
