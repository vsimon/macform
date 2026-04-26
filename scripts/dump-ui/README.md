How to run:

`osascript scripts/dump-ui/dump-ui.applescript "System Settings" "Displays"`

The format per line is:

```text
<class> [name] =value (description)
```

**`<class>`** — the accessibility class: group, checkbox, button, scroll area, static text, pop up button, slider, etc.

**`[name]`** — the element's accessible name, if set. For checkboxes and buttons this is usually the label the user sees. Omitted when missing.

**`=value`** — the element's current value. =0/=1 for checkboxes, the selected option for pop-ups, slider position, etc. Omitted when missing.

**(description)** — a secondary description when it differs from the class name. E.g. (switch) on a checkbox, (standard window) on a window, (Sidebar) on an outline. Omitted when it would just repeat the class.

**Indentation** = 2 spaces per depth level. The depth tells you exactly how to construct the AppleScript path: each indent level is one element reference you need to traverse. A group two levels below a scroll area is group N of scroll area M of ….