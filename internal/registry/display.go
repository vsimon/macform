package registry

import "fmt"

func displaySettingsScript(innerAction string) string {
	return fmt.Sprintf(`
if application "System Settings" is running then
	tell application "System Events"
		tell process "System Settings"
			try
				repeat while (count of sheets of window 1) > 0
					key code 53
					delay 0.3
				end repeat
			end try
		end tell
	end tell
end if
do shell script "open -g x-apple.systempreferences:com.apple.Displays-Settings.extension"
tell application "System Settings" to activate
delay 1
tell application "System Events"
	tell process "System Settings"
		set cb to missing value
		set waited to 0
		repeat while waited < 10
			try
				set cb to checkbox "Automatically adjust brightness" of group 1 of scroll area 2 of group 1 of group 3 of splitter group 1 of group 1 of window 1
				exit repeat
			end try
			delay 0.5
			set waited to waited + 1
		end repeat
		if cb is missing value then
			tell application "System Settings" to quit
			error "Automatically adjust brightness checkbox not found"
		end if
		%s cb
	end tell
end tell
tell application "System Settings" to quit
`, innerAction)
}

var autoBrightnessReadScript = displaySettingsScript("return value of")

// a click always reaches the desired state because Write only calls this when current != desired.
func autoBrightnessWriteScript(_ string) string {
	return displaySettingsScript("click")
}
