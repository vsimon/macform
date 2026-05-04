# macform — Product Requirements

## Executive Summary

`macform` is a Go CLI tool that lets users define macOS system settings in a YAML spec file and reconcile the system state against that spec — generating diffs, applying changes, and snapshotting current settings. It is designed for developers who commit their configuration file to a repo and want a reliable, repeatable way to bootstrap new machines.

## Problem Statement

When setting up a new Mac, developers spend significant time manually tweaking System Settings that they've already carefully tuned on their previous machine. There's no native tool to declaratively capture these settings and reproduce them. Existing solutions (shell scripts, Ansible) are either too brittle or require heavy infrastructure for a personal config use case.

## User Stories

### US-1: Bootstrap from spec file

**As a** developer setting up a new Mac,
**I want to** run `macform apply -f ~/dotfiles/macform.yaml`
**So that** my Dock and Finder settings are configured exactly as I like them within seconds.
**Acceptance criteria:**

- All keys defined in the spec are written to the system
- Keys with `null` values are deleted (restoring macOS defaults)
- Changes that require a process restart (Dock, Finder) trigger that restart automatically
- Command exits non-zero if any write fails

### US-2: Drift detection

**As a** developer who periodically tweaks settings manually,
**I want to** run `macform plan -f ~/dotfiles/macform.yaml`
**So that** I can see which settings have drifted from my committed spec before deciding whether to update the spec or revert the system.
**Acceptance criteria:**

- Header line: `macform will perform the following actions:`
- Only settings that will change or be deleted are shown — unchanged settings are omitted from output
- Changes are grouped by section (e.g. `# dock`, `# finder`) with individual settings indented beneath
- New settings (key does not exist on the system): `+ key = (not set) -> desired` in green
- Changed settings (key exists but differs): `~ key = current -> desired` in yellow
- Null-targeted settings (key will be deleted): `- key = current -> (deleted)` in red
- Summary line: `Plan: N to add, N to change, N to remove.`
- Zero exit code even when differences exist (it's a read-only operation)

### US-3: Snapshot current settings

**As a** developer setting up `macform` for the first time,
**I want to** run `macform generate -f macform.yaml`
**So that** I get a spec file pre-populated with my current system settings that I can commit to my dotfiles.
**Acceptance criteria:**

- Output file is valid YAML — bare keys and values only, no comment blocks
- All supported settings that are currently set appear with their current values (decoded to human-readable form per FR-9)
- Keys not currently set on the system are omitted entirely
- File is written atomically (temp file + rename)
- For the annotated reference with descriptions and valid values, see `examples/macform.yaml`

### US-4: Install via Homebrew

**As a** developer who wants to use macform on a new Mac,
**I want to** install it with a single `brew install` command
**So that** I don't need to manually download binaries or build from source.
**Acceptance criteria:**

- A Homebrew tap exists and is kept current with each release
- `brew install --cask vsimon/tap/macform` installs the correct binary for the running arch (arm64 or amd64)
- The installed binary passes `macform --version`

### US-5: Discover available settings

**As a** developer who is new to macform,
**I want to** open `examples/macform.yaml` and see every supported setting with documentation inline
**So that** I can copy what I want into my own spec without reading source code or external docs.
**Acceptance criteria:**

- Every supported setting is present with a sensible default value
- Each setting has a YAML comment block above it documenting: description, type, and valid values (for enums/strings)
- The file is valid YAML and loadable by `macform plan`
- The file is kept in sync with the registry — adding a setting to the registry requires adding it to `examples/macform.yaml`

## Functional Requirements

### FR-1: Spec file (YAML)

- Top-level keys: `dock`, `finder`, `display`, `battery`, `control-center`, `trackpad`, `keyboard`, `hot-corners`
- Each section contains setting key/value pairs
- Supported value types: `bool`, `int`, `float`, `string`
- `null` (or `~`) as a value means "delete this defaults key"
- Keys omitted from the spec are not touched

### FR-2: Spec file resolution

- Default filename: `macform.yaml`
- Resolution order: (1) `--file / -f` flag, (2) `macform.yaml` in the current working directory
- If no file is found and no flag is given, print a clear error: `No spec file found. Run 'macform generate' to create one, or pass --file.`

### FR-3: `plan` command

- Resolves and validates spec file per FR-2 and FR-6
- Reads current system state for all settings defined in the spec
- Prints diff grouped by section, following the output format defined in US-2
- If no changes: prints `No changes. System matches spec.`
- Exits 0 always (read-only)

### FR-4: `apply` command

- Resolves and validates spec file per FR-2 and FR-6
- Computes changes and **prints the full plan output** (same format as `plan` command) before prompting
- If no changes: prints `No changes. System matches spec.` and exits 0
- Prompts for user confirmation unless `--auto-approve` flag is set
  Confirmation prompt for apply must prompt the user for 'yes', use bold formatting in output:

```markdown
**Do you want to perform these actions?**
   macform will perform the actions described above.
   Only 'yes' will be accepted to approve.

   **Enter a value:**
```

- Applies each change via the appropriate provider; **stops at the first failure** and exits non-zero, reporting which setting failed and why
- After all writes succeed, restarts affected processes (Dock, Finder) once each
- Audit log of changes is printed to the console (in white) with a final summary (in bold green). As an example:

```markdown
  ~ dock
      ~ magnification: false → true
      $ killall Dock

**Apply complete! Resources: N added, N changed, N removed.**
```

### FR-5: `generate` command

- Reads all supported settings from the system
- Writes a **bare YAML snapshot** — keys and values only, no comment blocks (see `examples/macform.yaml` for the annotated reference)
- Omits settings that are not currently set on the system
- Output path: `--file / -f` flag if provided, otherwise `macform.yaml` in the current directory
- Warns if file already exists and prompts for confirmation (unless `--force`)

### FR-6: `--version` flag

- All commands inherit a `--version` / `-v` flag (or `macform version` subcommand)
- Output format: `macform v1.0.0 (commit: abc1234, built: 2026-04-19)`
- Values injected at build time via `-ldflags`; defaults to `dev` when built without ldflags (e.g. `go run .`)

### FR-7: `--no-color` flag

- Global flag `--no-color` disables ANSI color output for all commands
- Also respected when `NO_COLOR` environment variable is set (per [no-color.org](https://no-color.org) convention)

### FR-8: Spec validation

- The spec file is validated before any read or write operation
- Errors that abort immediately:
  - Invalid YAML syntax
  - Unknown top-level section (e.g. `dck:` — not a recognized section)
  - Unknown key within a section (e.g. `autohid: true` — likely a typo)
  - Wrong value type for a key (e.g. `tile-size: "big"` when an int is expected)
- Validation errors print the offending key/line and a suggestion where possible (e.g. `Unknown key "autohid" in dock — did you mean "autohide"?`)
- Warnings (non-fatal): unrecognized value for a string setting that has a known value map

### FR-9: Human-readable spec values (encode/decode)

- Spec values must use user-facing, human-readable names as they appear in macOS System Settings — not raw system values
- Where system values are opaque or abbreviated, the registry declares a bidirectional value map used to encode (spec → system) on write and decode (system → spec) on read
- Example: Finder's view style is stored as `icnv`, `Nlsv`, `clmv`, `Flwv` — the spec uses `icon`, `list`, `column`, `gallery`
- The annotated `examples/macform.yaml` lists valid values using the human-readable names only; raw system values are never exposed to the user
- `generate` decodes system values into spec values before writing the output file
- If a system value is encountered that has no mapping, it is written as-is with a warning

### FR-10: Settings provider abstraction

- A `Provider` interface abstracts how settings are read/written
- `defaults` provider: reads/writes via macOS `defaults` CLI; supports a `-currentHost` variant for ByHost-domain settings (e.g. control-center items that are machine-specific)
- `osascript` provider: GUI scripting for settings not accessible via `defaults` (e.g. display brightness, battery dim)
- Each setting in the registry declares which provider it uses

### FR-11: Settings support

**Dock** (`com.apple.dock`):

| Spec Key | Defaults Key | Type | Value Map |
| --- | --- | --- | --- |
| autohide | autohide | bool | — |
| tile-size | tilesize | int | — |
| orientation | orientation | string | — (`bottom`, `left`, `right` are already human-readable) |
| minimize-to-application | minimize-to-application | bool | — |
| show-recents | show-recents | bool | — |
| magnification | magnification | bool | — |
| large-size | largesize | int | — |
| min-effect | mineffect | string | — (`genie`, `scale`, `suck` match macOS UI labels) |
| scroll-to-open | scroll-to-open | bool | — |

**Finder** (mixed domains):

| Spec Key | Defaults Key | Domain | Type | Value Map |
| --- | --- | --- | --- | --- |
| show-hidden-files | AppleShowAllFiles | com.apple.finder | bool | — |
| show-extensions | AppleShowAllExtensions | NSGlobalDomain | bool | — |
| show-path-bar | ShowPathbar | com.apple.finder | bool | — |
| show-status-bar | ShowStatusBar | com.apple.finder | bool | — |
| default-view-style | FXPreferredViewStyle | com.apple.finder | string | `icon`→`icnv`, `list`→`Nlsv`, `column`→`clmv`, `gallery`→`Flwv` |
| warn-on-extension-change | FXEnableExtensionChangeWarning | com.apple.finder | bool | — |
| new-window-target | NewWindowTarget | com.apple.finder | string | `recents`→`PfAF`, `home`→`PfHm`, `desktop`→`PfDe`, `documents`→`PfDo`, `computer`→`PfCm`, `volumes`→`PfVo`, `icloud-drive`→`PfID` |

**Display** (osascript GUI scripting):

| Spec Key | Provider | Type | macOS Default |
| --- | --- | --- | --- |
| auto-brightness | osascript | bool | true |

**Battery** (osascript GUI scripting):

| Spec Key | Provider | Type | macOS Default |
| --- | --- | --- | --- |
| slightly-dim-on-battery | osascript | bool | true |

**Control Center** (`com.apple.controlcenter`, written with `-currentHost` where noted):

| Spec Key | Defaults Key | Type | Value Map | Notes |
| --- | --- | --- | --- | --- |
| show-battery | Battery | bool (int) | `true`→`18`, `false`→`24` | — |
| show-battery-percentage | BatteryShowPercentage | bool | — | — |
| show-bluetooth | Bluetooth | bool (int) | `true`→`18`, `false`→`24` | — |
| show-sound | Sound | string (int) | `always`→`18`, `when-active`→`2`, `never`→`24` | `-currentHost` |
| show-spotlight | MenuItemHidden | bool (int) | `true`→`0`, `false`→`1` | domain: `com.apple.Spotlight` |
| show-wifi | WiFi | bool (int) | `true`→`18`, `false`→`24` | — |

**Trackpad** (dual-domain: `com.apple.AppleMultitouchTrackpad` + `com.apple.driver.AppleBluetoothMultitouch.trackpad`):

| Spec Key | Defaults Key | Type | Value Map |
| --- | --- | --- | --- |
| tap-to-click | Clicking | bool | — |
| tracking-speed | com.apple.trackpad.scaling (NSGlobalDomain) | float | — |
| dragging-style | Dragging + DragLock + TrackpadThreeFingerDrag | string | `with-drag-lock`, `without-drag-lock`, `three-finger-drag` (compound write) |

**Keyboard** (`NSGlobalDomain` unless noted):

| Spec Key | Defaults Key | Type | Value Map |
| --- | --- | --- | --- |
| repeat-rate | KeyRepeat | int | — |
| repeat-delay | InitialKeyRepeat | int | — |
| function-keys | com.apple.keyboard.fnState | string | `special`→`0`, `standard`→`1` |
| function-key-action | AppleFnUsageType (`com.apple.HIToolbox`) | string | `do-nothing`→`0`, `change-input-source`→`1`, `show-emoji`→`2`, `start-dictation`→`3` |
| auto-capitalize | NSAutomaticCapitalizationEnabled | bool | — |
| auto-correct | NSAutomaticSpellingCorrectionEnabled | bool | — |

**Hot Corners** (`com.apple.dock`, triggers `killall Dock`):

| Spec Key | Defaults Key | Type | Value Map |
| --- | --- | --- | --- |
| top-left | wvous-tl-corner | string | `no-op`→`1`, `mission-control`→`2`, … |
| top-left-modifier | wvous-tl-modifier | string | `none`→`0`, `shift`→`131072`, … |
| top-right | wvous-tr-corner | string | (same map as corners) |
| top-right-modifier | wvous-tr-modifier | string | (same map as modifiers) |
| bottom-left | wvous-bl-corner | string | (same map as corners) |
| bottom-left-modifier | wvous-bl-modifier | string | (same map as modifiers) |
| bottom-right | wvous-br-corner | string | (same map as corners) |
| bottom-right-modifier | wvous-br-modifier | string | (same map as modifiers) |

### FR-12: Annotated example spec

- A file `examples/macform.yaml` is committed to the repo
- Contains every setting supported by the registry with a sensible default value
- Each setting is preceded by a YAML comment block with:
  - `# Description:` — what the setting does
  - `# Type:` — `bool`, `int`, `float`, or `string`
  - `# Valid values:` — enumerated options for string settings (e.g. `bottom | left | right`); omitted for free-form types
  - `# Default:` — the macOS system default if the key is deleted
- The file must be valid YAML parseable by `macform plan`
- Treated as a living document: every new setting added to the registry must have a corresponding entry here

Example structure:

```yaml
dock:
  # Description: Automatically hide and show the Dock
  # Type: bool
  autohide: true

  # Description: Size of Dock icons in pixels
  # Type: int
  tile-size: 48

  # Description: Position of the Dock on screen
  # Type: string
  # Valid values: bottom | left | right
  orientation: bottom
```

### FR-13: Build & release pipeline

- **Tooling**: `mise` manages all repo tools (Go, GoReleaser) via `mise.toml` at the repo root — no global installs required
- **CI**: GitHub Actions installs mise, then uses it to activate Go and GoReleaser; runs `go build`, `go vet`, `go test ./...` on every push
- **Release trigger**: Pushing a semver tag (`v*.*.*`) kicks off the release workflow
- **Tool**: GoReleaser handles cross-compilation, archive creation, checksum generation, and GitHub Release publication
- **Targets**: `darwin/amd64` (Intel) and `darwin/arm64` (Apple Silicon); no other platforms
- **Archives**: `.tar.gz` per platform, plus a `checksums.txt` with SHA-256 digests
- **Homebrew tap**: GoReleaser auto-updates `vsimon/homebrew-tap` with a generated formula after each release
- **Version embedding**: Build injects `version`, `commit`, and `date` via `-ldflags` so `macform --version` reports accurate release info
- **Commit convention**: All commits must follow [Conventional Commits](https://www.conventionalcommits.org/) (`type(scope): description`). This is enforced via a PR title check in CI. Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `build`, `ci`, `revert`

### FR-14: Exit codes

| Command | Condition | Exit code |
| --- | --- | --- |
| `plan` | Always | 0 |
| `apply` | No changes | 0 |
| `apply` | All changes applied successfully | 0 |
| `apply` | Write failure (stops at first) | 1 |
| `apply` | User aborts at confirmation prompt | 1 |
| `generate` | File written successfully | 0 |
| `generate` | File write error | 1 |
| `generate` | User aborts overwrite prompt | 1 |
| Any command | Spec file not found | 1 |
| Any command | Spec validation error | 1 |

## Non-Functional Requirements

- **Single binary**: No runtime dependencies. Distributed as a compiled Go binary.
- **Fast**: `plan` should complete in under 2 seconds for a full spec.
- **Safe**: `plan` is always read-only. `apply` prompts before writing unless `--auto-approve`.
- **macOS compatibility**: Target macOS 26 (Tahoe) to start. Compatibility with multiple versions is a future goal.
- **Extensible**: Adding a new provider requires only: (a) implement the `Provider` interface, (b) add entries to the registry with the new provider name.

## Success Criteria

1. `generate` produces a valid, loadable YAML spec from a real Mac.
2. `plan` correctly identifies zero differences when system matches spec.
3. `plan` correctly shows diffs when settings differ (tested manually).
4. `apply --auto-approve` changes all differing settings and restarts affected apps.
5. After `apply`, running `plan` again shows "No changes."
6. `go build` succeeds with no warnings.
7. Pushing a `v*.*.*` tag triggers CI, produces GitHub Release with arm64 + amd64 archives, and updates the Homebrew formula.
8. `brew install vsimon/tap/macform` installs a working binary on a clean Mac.

## Constraints & Assumptions

- **Go only**: No external scripting runtimes required.
- **mise**: All developer tooling (Go version, GoReleaser) is declared in `mise.toml`. Contributors run `mise install` to get the exact versions; CI does the same.
- **YAML spec**: Chosen for null-value support and familiarity.
- **`defaults` CLI**: Used for all initial settings. Assumes macOS ships `defaults` (it always has).
- **Single user**: Manages settings for the current user only (no `sudo` required for initial settings).
- **null = delete**: A YAML `null` / `~` value means `defaults delete <domain> <key>`.

## Out of Scope

- GUI or menu bar application
- Third-party app settings (VS Code, JetBrains, etc.)
- Fetching spec files from remote URLs or shared repositories
- Windows or Linux support
- Rollback / undo after apply
