# TODO

## project setup
- [x] **T001 · Chore · P2: initialize go module and directory structure**
    - **Context:** Detailed Build Step 1
    - **Action:**
        1. Run `go mod init` with module path.
        2. Create root folders: `cmd/ward`, `internal/{config,core,git,llm,log,shim,ci,errors,types,util}`.
    - **Done-when:**
        1. `go.mod` exists.
        2. All directories compile (no missing package errors).
    - **Verification:**
        1. Run `go build ./...` without errors.
    - **Depends-on:** none

- [x] **T037 · Chore · P2: Verify current golangci-lint version**
    - **Context:** Step 1 from `CONSULTANT-PLAN.md`. The current `golangci-lint` version (reported as v2.1.1) predates the modern configuration schema (v2) and is causing `unsupported version` errors. We need to confirm the installed version before updating.
    - **Action:** Open a terminal in the project directory (`/Users/phaedrus/Development/ward/`) and run `golangci-lint --version`. Record the output.
    - **Done-when:** The output of the command, showing the installed `golangci-lint` version, is known and recorded.
    - **Depends-on:** T001

- [x] **T038 · Chore · P2: Update golangci-lint to the latest stable version**
    - **Context:** Step 1 from `CONSULTANT-PLAN.md`. The current version is too old. Updating to a recent stable version (>= v1.50.0) is necessary for compatibility with the required configuration schema and to benefit from bug fixes and new linters. The plan recommends using `go install`.
    - **Action:** Execute the command `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` to install or update `golangci-lint` to the latest stable version available via Go's tooling. Alternatively, use the appropriate package manager command if installed differently (e.g., `brew upgrade golangci-lint`).
    - **Done-when:** The update command completes successfully.
    - **Depends-on:** T037

- [x] **T039 · Chore · P2: Verify the updated golangci-lint version**
    - **Context:** Step 1 from `CONSULTANT-PLAN.md`. Confirm that the update performed in T038 was successful and the `golangci-lint` command now executes the newly installed version.
    - **Action:** Run `golangci-lint --version` again in the project directory.
    - **Done-when:** The output confirms a recent version of `golangci-lint` (e.g., v1.5x.y or newer) is now installed and accessible via the command line.
    - **Depends-on:** T038

- [x] **T040 · Chore · P2: Create the initial .golangci.yml configuration file**
    - **Context:** Step 2 from `CONSULTANT-PLAN.md`. A configuration file (`.golangci.yml`) using the `version: "2"` schema is required to define the linters and settings for the project, replacing the incompatible older format.
    - **Action:** Create a file named `.golangci.yml` in the project root directory (`/Users/phaedrus/Development/ward/`). Paste the complete YAML content provided in `CONSULTANT-PLAN.md` (starting with `version: "2"`) into this file.
    - **Done-when:** The `.golangci.yml` file exists in the project root, contains the base configuration from the plan, and is saved.
    - **Depends-on:** T039

- [ ] **T041 · Chore · P2: Customize .golangci.yml with project-specific details**
    - **Context:** Step 2 from `CONSULTANT-PLAN.md`. The base `.golangci.yml` configuration contains placeholders that must be tailored to the specific project (`ward`). Specifically, the Go version and the local module path prefix need to be set correctly.
    - **Action:** Edit the `.golangci.yml` file:
        1. Update the `go:` value under the `run:` section to match the Go version specified in the project's `go.mod` file (e.g., `go: '1.22'`).
        2. Update the `local-prefixes:` value under `linters-settings.goimports:` to match the project's module path found in `go.mod` (e.g., `local-prefixes: github.com/phaedrus-io/ward`).
        3. Add and commit the `.golangci.yml` file to version control.
    - **Done-when:** The `go:` version and `local-prefixes:` in `.golangci.yml` accurately reflect the project's `go.mod` file, and the configuration file is committed to the repository.
    - **Depends-on:** T040

- [ ] **T042 · Chore · P2: Run golangci-lint to validate configuration and initial state**
    - **Context:** Step 3 from `CONSULTANT-PLAN.md`. After creating and customizing the configuration, verify that the updated `golangci-lint` tool can parse the `.golangci.yml` file correctly and execute without configuration errors on the existing codebase.
    - **Action:** Navigate to the project root directory (`/Users/phaedrus/Development/ward/`) in the terminal and run `golangci-lint run ./...`.
    - **Done-when:** The command completes successfully (exit code 0) without reporting any configuration parsing errors (the `unsupported version` error must be resolved). Ideally, it reports no linting violations given the strict config and likely minimal codebase, but the primary goal here is config validation.
    - **Depends-on:** T041

- [ ] **T043 · Documentation · P3: Add linting instructions to README.md**
    - **Context:** Step 4 from `CONSULTANT-PLAN.md` and the project's development philosophy emphasize documenting tooling setup. This ensures developers know how to install and run the linter locally.
    - **Action:** Edit the project's `README.md` file. Add a "Linting" subsection under a "Development" or similar section, including the markdown content provided in Step 4 of `CONSULTANT-PLAN.md`. Ensure it clearly explains the tool (`golangci-lint`), the configuration file (`.golangci.yml`), how to install/check the version, and the command to run it (`golangci-lint run ./...`). Commit the changes.
    - **Done-when:** `README.md` contains a clear and accurate "Linting" section detailing the setup and local execution instructions, and the changes are committed.
    - **Depends-on:** T042

- [ ] **T044 · Meta · P3: Mark original task T002 as completed**
    - **Context:** The work originally scoped in task T002 (`configure golangci-lint`) has been fully decomposed into and addressed by tasks T037 through T043 based on the detailed `CONSULTANT-PLAN.md`.
    - **Action:** Edit the `TODO.md` file. Locate the line for task `T002 · Chore · P2: configure golangci-lint`. Change its status marker from `[ ]` (or `[~]`) to `[x]`. Commit the updated `TODO.md`.
    - **Done-when:** Task T002 in `TODO.md` is marked as completed `[x]`.
    - **Depends-on:** T043

- [ ] **T003 · Chore · P2: create Makefile for common tasks**
    - **Context:** Detailed Build Step 1
    - **Action:**
        1. Add `Makefile` with targets: `lint`, `test`, `build`.
        2. Wire commands: `golangci-lint run`, `go test ./...`, `go build ./cmd/ward`.
    - **Done-when:**
        1. `make lint`, `make test`, `make build` succeed.
    - **Depends-on:** T001

## internal/types
- [ ] **T004 · Chore · P2: define shared types**
    - **Context:** Detailed Build Step 2
    - **Action:**
        1. Create `internal/types` package.
        2. Define `ReviewRequest`, `ReviewResult`, `LogEntry`, `Config` structs.
    - **Done-when:**
        1. Types compile without errors.
    - **Depends-on:** T001

## internal/errors
- [ ] **T005 · Chore · P2: define sentinel errors**
    - **Context:** Detailed Build Step 3
    - **Action:**
        1. Create `internal/errors` package.
        2. Declare `ErrCheckFailed`, `ErrLogUpdateNotFound`, `ErrExternalCmdFailed`, `ErrConfigInvalid` errors.
    - **Done-when:**
        1. Errors compile and are addressable from other packages.
    - **Depends-on:** T001

## internal/git
- [ ] **T006 · Chore · P2: define git.Client interface**
    - **Context:** Detailed Build Step 4
    - **Action:**
        1. Create `internal/git/client.go`.
        2. Declare methods: `GetStagedDiff`, `GetCommitMessage`, `GetCurrentBranch`, `GetHeadCommitHash`, `GetGitDir`, `WritePlaceholderHash`, `ReadPlaceholderHash`, `CleanPlaceholderHash`.
    - **Done-when:**
        1. Interface compiles.
    - **Depends-on:** T005

- [ ] **T018 · Feature · P2: implement git.Client adapter using util exec wrapper**
    - **Context:** Detailed Build Step 8
    - **Action:**
        1. Use `internal/util` exec wrapper to call Git CLI.
        2. Implement placeholder file I/O in `.git/ward_placeholder`.
    - **Done-when:**
        1. Methods return expected outputs.
        2. Error handling wraps errors with `ErrExternalCmdFailed`.
    - **Depends-on:** T006, T017

- [ ] **T019 · Test · P2: add integration tests for git.Client**
    - **Context:** Detailed Build Step 8
    - **Action:**
        1. Set up temporary Git repo.
        2. Test all `git.Client` methods under normal and edge-case states.
    - **Done-when:**
        1. Integration tests pass.
    - **Depends-on:** T018

## internal/llm
- [ ] **T007 · Chore · P2: define llm.Client interface**
    - **Context:** Detailed Build Step 4
    - **Action:**
        1. Create `internal/llm/client.go`.
        2. Declare `AnalyzeChanges(ctx, req) (*types.ReviewResult, error)`.
    - **Done-when:**
        1. Interface compiles.
    - **Depends-on:** T006

- [ ] **T020 · Feature · P2: implement llm.Client adapter using util exec wrapper**
    - **Context:** Detailed Build Step 9
    - **Action:**
        1. Call `claude` CLI via `internal/util` wrapper.
        2. Parse LLM output into `ReviewResult`, handle PASS/WARN/FAIL.
    - **Done-when:**
        1. Adapter returns correct `ReviewResult` or error.
    - **Depends-on:** T007, T017

- [ ] **T021 · Test · P2: add integration tests for llm.Client**
    - **Context:** Detailed Build Step 9
    - **Action:**
        1. Mock `os/exec` to simulate `claude` output.
        2. Verify parsing for each status.
    - **Done-when:**
        1. Tests cover PASS, WARN, FAIL, parse errors.
    - **Depends-on:** T020

## internal/log
- [ ] **T008 · Chore · P2: define log.Writer interface**
    - **Context:** Detailed Build Step 4
    - **Action:**
        1. Create `internal/log/writer.go`.
        2. Declare `WriteEntry`, `UpdateEntryCommitHash`.
    - **Done-when:**
        1. Interface compiles.
    - **Depends-on:** T005

- [ ] **T022 · Feature · P2: implement log.Writer adapter with JSON Lines and locking**
    - **Context:** Detailed Build Step 10
    - **Action:**
        1. Append entries as JSON Lines to `.ward-warnings.log`.
        2. Implement atomic update and file lock for placeholder updates.
    - **Done-when:**
        1. Log file format matches spec.
        2. `UpdateEntryCommitHash` handles placeholder correctly.
    - **Depends-on:** T008, T017

- [ ] **T023 · Test · P2: add integration tests for log.Writer**
    - **Context:** Detailed Build Step 10
    - **Action:**
        1. Use temp directory to write and update log entries.
        2. Verify locking prevents races.
    - **Done-when:**
        1. Integration tests pass without flakiness.
    - **Depends-on:** T022

## internal/shim
- [ ] **T009 · Chore · P2: define shim.Generator interface**
    - **Context:** Detailed Build Step 4
    - **Action:**
        1. Create `internal/shim/generator.go`.
        2. Declare `Generate(ctx, manager, hookType) (string, error)`.
    - **Done-when:**
        1. Interface compiles.
    - **Depends-on:** T005

- [ ] **T024 · Feature · P2: implement shim.Generator for pre-commit**
    - **Context:** Detailed Build Step 11
    - **Action:**
        1. Generate hook script calling `ward check` or `ward log`.
        2. Support at least `pre-commit` manager.
    - **Done-when:**
        1. Script content matches expected template.
    - **Depends-on:** T009

- [ ] **T025 · Test · P2: add unit tests for shim.Generator**
    - **Context:** Detailed Build Step 11
    - **Action:**
        1. Verify `Generate` output for known manager and hook types.
    - **Done-when:**
        1. Tests cover valid and invalid manager.
    - **Depends-on:** T024

## internal/ci
- [ ] **T010 · Chore · P2: define ci.Detector interface**
    - **Context:** Detailed Build Step 4
    - **Action:**
        1. Create `internal/ci/detector.go`.
        2. Declare `IsCI() bool`.
    - **Done-when:**
        1. Interface compiles.
    - **Depends-on:** T005

- [ ] **T014 · Feature · P2: implement ci.Detector**
    - **Context:** Detailed Build Step 6
    - **Action:**
        1. Detect CI via `CI` or `GITHUB_ACTIONS` env vars.
    - **Done-when:**
        1. `IsCI()` returns correct values.
    - **Depends-on:** T010

- [ ] **T015 · Test · P2: add unit tests for ci.Detector**
    - **Context:** Detailed Build Step 6
    - **Action:**
        1. Set/unset env vars to exercise both branches.
    - **Done-when:**
        1. Unit tests pass.
    - **Depends-on:** T014

## internal/config
- [ ] **T012 · Feature · P2: implement config loader**
    - **Context:** Detailed Build Step 5
    - **Action:**
        1. Read env vars into `Config` struct.
        2. Validate required fields.
    - **Done-when:**
        1. Loader fails on invalid/missing vars.
    - **Depends-on:** T004, T005

- [ ] **T013 · Test · P2: add unit tests for config loader**
    - **Context:** Detailed Build Step 5
    - **Action:**
        1. Mock env for valid and invalid cases.
    - **Done-when:**
        1. Tests cover success and error paths.
    - **Depends-on:** T012

## internal/util
- [ ] **T016 · Feature · P2: implement util package**
    - **Context:** Detailed Build Step 7
    - **Action:**
        1. Add UUID correlation ID generator.
        2. Build `exec.CommandContext` wrapper without shell interpolation.
    - **Done-when:**
        1. Functions available and documented.
    - **Depends-on:** T005

- [ ] **T017 · Test · P2: add unit tests for util**
    - **Context:** Detailed Build Step 7
    - **Action:**
        1. Test correlation ID format.
        2. Test exec wrapper handles exit codes and stderr.
    - **Done-when:**
        1. Unit tests pass.
    - **Depends-on:** T016

## internal/core
- [ ] **T011 · Chore · P2: define core interfaces**
    - **Context:** Detailed Build Step 4
    - **Action:**
        1. Create `internal/core/interfaces.go`.
        2. Declare `Checker`, `LogUpdater`, `Initializer` interfaces.
    - **Done-when:**
        1. Interfaces compile.
    - **Depends-on:** T005, T006, T007, T008, T009, T010

- [ ] **T026 · Feature · P2: implement core.Checker logic**
    - **Context:** Detailed Build Step 12
    - **Action:**
        1. Inject `GitClient`, `LLMClient`, `log.Writer`, `ci.Detector`.
        2. Implement empty-diff and CI skip flows, placeholder write, LLM call, entry write.
    - **Done-when:**
        1. Checker returns correct exit codes.
    - **Depends-on:** T011, T018, T020, T022, T014

- [ ] **T027 · Feature · P2: implement core.LogUpdater logic**
    - **Context:** Detailed Build Step 12
    - **Action:**
        1. Inject `GitClient`, `log.Writer`, `ci.Detector`.
        2. Implement placeholder read, actual hash update, entry update.
    - **Done-when:**
        1. LogUpdater updates log file and cleans placeholder.
    - **Depends-on:** T011, T018, T022, T014

- [ ] **T028 · Feature · P2: implement core.Initializer logic**
    - **Context:** Detailed Build Step 12
    - **Action:**
        1. Inject `shim.Generator`, `config.Loader`.
        2. Generate and write hook scripts for specified manager.
    - **Done-when:**
        1. Scripts created at correct paths.
    - **Depends-on:** T011, T024, T012

- [ ] **T029 · Test · P2: add unit tests for internal/core**
    - **Context:** Detailed Build Step 12
    - **Action:**
        1. Use mocks for all dependencies.
        2. Cover success, error, and edge-case flows.
    - **Done-when:**
        1. Code coverage > 95% for core package.
    - **Depends-on:** T026, T027, T028

## cmd/ward
- [ ] **T030 · Feature · P2: implement Cobra CLI commands**
    - **Context:** Detailed Build Step 13
    - **Action:**
        1. Add `cmd/ward` with subcommands `check`, `log`, `init`.
        2. Wire dependency injection and exit codes.
    - **Done-when:**
        1. `ward --help` shows commands.
        2. Commands return correct codes for mock core.
    - **Depends-on:** T012, T011

- [ ] **T031 · Test · P2: add unit tests for cmd/ward**
    - **Context:** Detailed Build Step 13
    - **Action:**
        1. Test flag parsing and command dispatch.
        2. Simulate errors to verify exit codes.
    - **Done-when:**
        1. Tests cover happy and error paths.
    - **Depends-on:** T030

## e2e
- [ ] **T032 · Test · P2: create E2E test script for git workflow**
    - **Context:** Detailed Build Step 14
    - **Action:**
        1. Script simulating `git add`, `git commit` with `ward` hook.
        2. Verify exit codes and `.ward-warnings.log` contents.
    - **Done-when:**
        1. E2E tests pass on CI.
    - **Depends-on:** T030, T019, T021, T023

## release
- [ ] **T033 · Chore · P2: configure GoReleaser**
    - **Context:** Detailed Build Step 15
    - **Action:**
        1. Add `.goreleaser.yml` for cross-platform static builds.
    - **Done-when:**
        1. Running `goreleaser --snapshot --skip-publish` succeeds.
    - **Depends-on:** T030

## documentation
- [ ] **T034 · Chore · P2: update README.md and docs**
    - **Context:** Detailed Build Step 16
    - **Action:**
        1. Document installation, usage, config, log format.
        2. Embed updated architecture diagram.
    - **Done-when:**
        1. README renders correctly.
    - **Depends-on:** T001, T030

- [ ] **T035 · Chore · P2: add Go doc comments for exports**
    - **Context:** All packages
    - **Action:**
        1. Add `//` comments explaining purpose of public types/functions.
    - **Done-when:**
        1. `godoc` shows descriptions.
    - **Depends-on:** T004–T031

## ci
- [ ] **T036 · Chore · P2: configure CI workflow**
    - **Context:** Testing Strategy & Automation Section
    - **Action:**
        1. Create CI pipeline to run `make lint`, `make test`, coverage check, `govulncheck`, and `make build`.
    - **Done-when:**
        1. Pull requests enforce all stages.
    - **Depends-on:** T002, T003, T032

### Clarifications & Assumptions
- [ ] **Issue:** confirm behavior if `.ward-warnings.log` is missing during `ward log` update
    - **Context:** Open Questions #1
    - **Blocking?:** yes
- [ ] **Issue:** decide if PASS results should be logged to `.ward-warnings.log`
    - **Context:** Open Questions #2
    - **Blocking?:** yes
- [ ] **Issue:** finalize initial list of supported hook managers for `ward init`
    - **Context:** Open Questions #3
    - **Blocking?:** yes
- [ ] **Issue:** determine long-term plan for direct Claude API versus CLI usage
    - **Context:** Open Questions #4
    - **Blocking?:** no
- [ ] **Issue:** clarify constraints on content/format of development philosophy summary passed to LLM
    - **Context:** Open Questions #5
    - **Blocking?:** no
- [ ] **Issue:** define size limits or content types to exclude from LLM review
    - **Context:** Open Questions #6
    - **Blocking?:** no
