# Shield CLI Demo - SSH
# Usage: VHS_CHROME_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" vhs demo-ssh.tape

Output docs/demo/demo-ssh.gif

Set FontSize 16
Set Width 1200
Set Height 600
Set Theme "Catppuccin Mocha"
Set TypingSpeed 80ms
Set Padding 20

# Set password silently, then clear screen before recording visible content
Hide
Type "export SHIELD_PASS='fengyi1545'"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 500ms
Show

# SSH connection (password shown as $SHIELD_PASS, not plaintext)
Type "shield ssh 172.16.3.137 --username root --auth-pass $SHIELD_PASS"
Sleep 1s
Enter
Sleep 15s
