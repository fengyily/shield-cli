# Shield CLI Demo - RDP
# Usage: VHS_CHROME_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" vhs demo-rdp.tape

Output docs/demo/demo-rdp.gif

Set FontSize 16
Set Width 1200
Set Height 600
Set Theme "Catppuccin Mocha"
Set TypingSpeed 80ms
Set Padding 20

# Set password silently, then clear screen before recording visible content
Hide
Type "export SHIELD_PASS='123456'"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 500ms
Show

# RDP connection (password shown as $SHIELD_PASS, not plaintext)
Type "shield rdp 172.16.3.28 --username REMOTE --auth-pass $SHIELD_PASS"
Sleep 1s
Enter
Sleep 15s
