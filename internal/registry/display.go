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
delay 1
tell application "System Events"
	tell process "System Settings"
		set cb to missing value
		set waited to 0
		repeat while waited < 10
			try
				set g3 to group 1 of group 3 of splitter group 1 of group 1 of window "Displays"
				-- Find the scroll area describing Built-in Display settings
				set settingsSA to missing value
				repeat with saIdx from 1 to count scroll areas of g3
					if description of scroll area saIdx of g3 contains "Built-in" then
						set settingsSA to scroll area saIdx of g3
						exit repeat
					end if
				end repeat
				-- Not visible: click each display thumbnail until Built-in settings appear
				if settingsSA is missing value then
					set thumbSA to missing value
					repeat with saIdx from 1 to count scroll areas of g3
						if description of scroll area saIdx of g3 is "Displays" then
							set thumbSA to scroll area saIdx of g3
							exit repeat
						end if
					end repeat
					if thumbSA is not missing value then
						repeat with btnIdx from 1 to count buttons of thumbSA
							click button btnIdx of thumbSA
							delay 0.7
							set g3 to group 1 of group 3 of splitter group 1 of group 1 of window "Displays"
							repeat with saIdx from 1 to count scroll areas of g3
								if description of scroll area saIdx of g3 contains "Built-in" then
									set settingsSA to scroll area saIdx of g3
									exit repeat
								end if
							end repeat
							if settingsSA is not missing value then exit repeat
						end repeat
					end if
				end if
				-- Search all groups in the settings scroll area for the checkbox
				if settingsSA is not missing value then
					repeat with gi from 1 to count groups of settingsSA
						try
							set cb to checkbox "Automatically adjust brightness" of group gi of settingsSA
							exit repeat
						end try
					end repeat
				end if
				if cb is not missing value then exit repeat
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
