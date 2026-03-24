# Shield CLI Demo - MySQL Plugin
# Usage: vhs demo-mysql.tape

Output docs/demo/demo-mysql.gif

Set FontSize 16
Set Width 1200
Set Height 600
Set Theme "Catppuccin Mocha"
Set TypingSpeed 80ms
Set Padding 20

# Set credentials silently, then clear screen before recording visible content
Hide
Type "export DB_PASS='root'"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 500ms
Show

# MySQL plugin connection (password shown as $DB_PASS, not plaintext)
Type "shield mysql 127.0.0.1:3306 --db-user root --db-pass $DB_PASS --db-name mydb"
Sleep 1s
Enter
Sleep 15s
