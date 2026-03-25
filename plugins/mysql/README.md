cd /Users/f1/Documents/fengyi/sheild-cli/plugins/mysql

# 先编译
go build -o shield-plugin-mysql .

# 启动插件，通过 echo 发送 start 指令（改成你的 MySQL 连接信息）
echo '{"action":"start","config":{"host":"127.0.0.1","port":3306,"user":"root","pass":"dataspace123","database":"","readonly":false}}' | ./shield-plugin-mysql

(echo '{"action":"start","config":{"host":"127.0.0.1","port":3306,"user":"root","pass":"dataspace123","database":"","readonly":false}}'; cat) | ./shield-plugin-mysql
