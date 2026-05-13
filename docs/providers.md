# Providers

Each setting in the registry is backed by a `Provider`.

## Available providers

### `Defaults`

Reads and writes via the macOS `defaults` CLI. The most common provider. Covers any setting stored in a macOS preferences domain (plist-backed). Prefer this over all other providers when a `defaults` key exists for the setting.

### `CurrentHostDefaults`

Same as `defaults` but passes `-currentHost` to scope reads and writes to the current machine. Used for settings stored per-host rather than per-user (e.g. some Control Center items).

### `MultiDefaults`

Reads from a primary domain; writes and deletes to the primary domain and one or more extra domains. Used when macOS requires the same key to be set in multiple preference domains (e.g. trackpad settings apply to both the built-in and Bluetooth trackpad domains).

### `Osascript`

Uses AppleScript GUI scripting via `osascript` to read and write settings that are not exposed through `defaults`. Requires System Events accessibility permissions. Use only as a last resort as GUI scripting is sensitive to macOS UI changes and more likely to break across OS updates.

### `DockAppPresence`

Reads and manipulates the Dock's `com.apple.dock.plist` directly via `/usr/libexec/PlistBuddy`. Used for managing which app entries are pinned to the Dock that `defaults` cannot address atomically.

### `DraggingStyle`

A compound custom provider that encodes the trackpad dragging style across three boolean keys (`Dragging`, `DragLock`, `TrackpadThreeFingerDrag`) written to both trackpad domains. There is no single `defaults` key for this concept; this provider synthesizes the state from the combination.

---

## Supported Settings

| Category | Setting | Description | Provider |
|----------|---------|-------------|----------|
| **appearance** | show-scroll-bars | When to show scroll bars | Defaults |
| **dock** | animate-opening-applications | Animate opening applications in the Dock | Defaults |
| | autohide | Auto-hide the Dock when the pointer is not near it | Defaults |
| | large-size | Maximum icon size in pixels when magnification is active | Defaults |
| | magnification | Enable Dock icon magnification when the pointer hovers over it | Defaults |
| | min-effect | Animation effect used when minimizing windows | Defaults |
| | minimize-to-application | Minimize windows into their application icon instead of the Dock | Defaults |
| | orientation | Position of the Dock on the screen | Defaults |
| | remove-apps | Bundle IDs of applications to remove from the Dock | DockAppPresence |
| | scroll-to-open | Scroll up on a Dock icon to show all opened windows across Spaces | Defaults |
| | show-recents | Show recently opened applications in the Dock | Defaults |
| | tile-size | Icon size in pixels for items in the Dock | Defaults |
| **finder** | auto-open-ro-volumes | Automatically open a Finder window when a read-only disk image is mounted | Defaults |
| | auto-open-rw-volumes | Automatically open a Finder window when a read-write disk image is mounted | Defaults |
| | default-view-style | Default view style for new Finder windows | Defaults |
| | desktop-show-external-drives | Show external hard drives on the desktop | Defaults |
| | desktop-show-hard-drives | Show internal hard drives on the desktop | Defaults |
| | desktop-show-removable-media | Show removable media (USB drives, DVDs) on the desktop | Defaults |
| | desktop-show-servers | Show mounted network servers on the desktop | Defaults |
| | new-window-target | Target folder for new Finder windows | Defaults |
| | no-ds-store-on-network | Prevent Finder from creating .DS_Store files on network volumes | Defaults |
| | no-ds-store-on-usb | Prevent Finder from creating .DS_Store files on USB volumes | Defaults |
| | quick-look-text-selection | Allow selecting and copying text in Quick Look previews | Defaults |
| | quit-menu-item | Allow quitting Finder via Cmd+Q, adding Quit to Finder's menu | Defaults |
| | show-extensions | Always show file name extensions in the Finder | Defaults |
| | show-hidden-files | Show hidden files and folders (those whose names start with a dot) | Defaults |
| | show-path-bar | Show the path bar at the bottom of Finder windows | Defaults |
| | show-posix-path-in-title | Show the full POSIX path to the current folder in the Finder title bar | Defaults |
| | show-status-bar | Show the status bar at the bottom of Finder windows | Defaults |
| | sort-folders-first | Keep folders on top when sorting by name in Finder windows | Defaults |
| | warn-on-extension-change | Show a warning dialog when changing a file's extension | Defaults |
| **display** | auto-brightness | Automatically adjust display brightness based on ambient light | Osascript |
| **battery** | slightly-dim-on-battery | Slightly dim the display when on battery power | Osascript |
| **control-center** | clock-show-day-of-week | Show the day of the week in the menu bar clock | Defaults |
| | clock-show-seconds | Show seconds in the menu bar clock | Defaults |
| | show-battery | Show the battery icon in the menu bar | Defaults |
| | show-battery-percentage | Show battery percentage in the menu bar | Defaults |
| | show-bluetooth | Show the Bluetooth icon in the menu bar | Defaults |
| | show-sound | Show the Sound icon in the menu bar | CurrentHostDefaults |
| | show-spotlight | Show the Spotlight icon in the menu bar | Defaults |
| | show-wifi | Show the Wi-Fi icon in the menu bar | Defaults |
| **trackpad** | dragging-style | Dragging style; also enables "Use trackpad for dragging" when set | DraggingStyle |
| | natural-scrolling | Scroll direction — content moves with fingers (natural) or opposite (traditional) | Defaults |
| | tap-to-click | Enable tap to click on the trackpad | MultiDefaults |
| | tracking-speed | Trackpad tracking speed | Defaults |
| **keyboard** | auto-capitalize | Automatically capitalize the first letter of sentences | Defaults |
| | auto-correct | Automatically correct spelling errors | Defaults |
| | double-space-period | Insert a period when the space bar is pressed twice | Defaults |
| | function-key-action | Action performed when pressing the Fn / Globe key | Defaults |
| | function-keys | Behavior of the function keys (F1, F2, etc.) | Defaults |
| | keyboard-navigation | Move focus between controls using Tab and Shift-Tab | Defaults |
| | press-and-hold | Enable press-and-hold for keys to show accented character options | Defaults |
| | repeat-delay | Delay until key repeat starts | Defaults |
| | repeat-rate | Key repeat rate | Defaults |
| | smart-dashes | Automatically replace double hyphens and other dash sequences with em/en dashes | Defaults |
| **hot-corners** | bottom-left | Action triggered when the cursor reaches the bottom-left corner | Defaults |
| | bottom-left-modifier | Modifier key required to trigger the bottom-left corner action | Defaults |
| | bottom-right | Action triggered when the cursor reaches the bottom-right corner | Defaults |
| | bottom-right-modifier | Modifier key required to trigger the bottom-right corner action | Defaults |
| | top-left | Action triggered when the cursor reaches the top-left corner | Defaults |
| | top-left-modifier | Modifier key required to trigger the top-left corner action | Defaults |
| | top-right | Action triggered when the cursor reaches the top-right corner | Defaults |
| | top-right-modifier | Modifier key required to trigger the top-right corner action | Defaults |
