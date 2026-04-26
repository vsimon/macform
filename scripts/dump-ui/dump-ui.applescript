#!/usr/bin/osascript
--
-- dump-ui.applescript — recursively dump the accessibility tree of a window.
--
-- Usage:
--   osascript scripts/dump-ui/dump-ui.applescript
--   osascript scripts/dump-ui/dump-ui.applescript "System Settings"
--   osascript scripts/dump-ui/dump-ui.applescript "System Settings" "Displays"
--
-- Pipe through grep to filter, e.g.:
--   osascript scripts/dump-ui/dump-ui.applescript "System Settings" "Displays" | grep -i brightness

on run argv
	set appName to "System Settings"
	set winTitle to ""
	if (count of argv) >= 1 then set appName to item 1 of argv
	if (count of argv) >= 2 then set winTitle to item 2 of argv

	tell application "System Events"
		tell process appName
			if winTitle is "" then
				return my dumpEl(window 1, 0)
			else
				return my dumpEl(window winTitle, 0)
			end if
		end tell
	end tell
end run

on dumpEl(el, depth)
	set ind to ""
	repeat depth times
		set ind to ind & "  "
	end repeat

	set ln to ind
	set kids to {}

	tell application "System Events"
		try
			set elCls to class of el
			set ln to ln & (elCls as string)
		end try
		try
			set nm to name of el
			if nm is not missing value and nm is not "" then
				set ln to ln & " [" & nm & "]"
			end if
		end try
		try
			set v to value of el
			if v is not missing value then
				set vs to v as string
				if vs is not "missing value" then set ln to ln & " =" & vs
			end if
		end try
		try
			set d to description of el
			if d is not missing value and d is not "" and d is not (class of el as string) then
				set ln to ln & " (" & d & ")"
			end if
		end try
		try
			set kids to every UI element of el
		end try
	end tell

	set txt to ln & linefeed

	repeat with kid in kids
		set txt to txt & my dumpEl(kid, depth + 1)
	end repeat

	return txt
end dumpEl
