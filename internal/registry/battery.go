package registry

import "fmt"

func batteryOptionsScript(postCheckboxAction string) string {
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
do shell script "open -g x-apple.systempreferences:com.apple.Battery-Settings.extension"
delay 1
tell application "System Events"
	tell process "System Settings"
		set btn to missing value
		set waited to 0
		repeat while waited < 10
			try
				set btn to button 1 of scroll area 1 of group 1 of group 3 of splitter group 1 of group 1 of window 1
				exit repeat
			end try
			delay 0.5
			set waited to waited + 1
		end repeat
		if btn is missing value then
			tell application "System Settings" to quit
			error "Options button not found in Battery settings"
		end if
		click btn
		delay 0.5
		set cb to missing value
		set waited to 0
		repeat while waited < 10
			try
				set cb to checkbox "Slightly dim the display on battery" of group 1 of scroll area 1 of group 1 of sheet 1 of window 1
				exit repeat
			end try
			delay 0.5
			set waited to waited + 1
		end repeat
		if cb is missing value then
			key code 53
			tell application "System Settings" to quit
			error "Slightly dim the display on battery checkbox not found"
		end if
		%s
	end tell
end tell
tell application "System Settings" to quit
`, postCheckboxAction)
}

var slightlyDimReadScript = batteryOptionsScript(`set cbVal to value of cb
		click button 1 of group 1 of sheet 1 of window 1
		return cbVal`)

func slightlyDimWriteScript(_ string) string {
	return batteryOptionsScript(`click cb
		click button 1 of group 1 of sheet 1 of window 1`)
}
