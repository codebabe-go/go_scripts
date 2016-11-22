# gitgo

## 功能介绍
* 一键提交代码, 不用重复做add, commit, pull, push等机械工作
* 直接用IDE...逼格不高

## 使用方法
  1. 配置go环境
     * OS X: brew install go
     * ubuntu: apt-get install go
     * others: https://golang.org/dl/ 选择合适版本安装
  2. 设置go path
     * Windows: 环境变量中配置
     * *nix: 配置.bash_profile $PATH:$GO_PATH
  3. 生成二进制文件
     go build ${project_path}/contrib/go_scripts/gitgo.go
  4. 运行
     假设你go build 时候的路径为build_path
     * Windows: cmd -> cd build_path -> gitgo.exe comment push_branch(待补充)
     * *nix: .bash_profile "alias gitgo=build_path/gitgo", 直接使用gitgo comment push_branch
  
## 缺点
  这里不去获取本地的user信息, 需要保障的前提是本地git命令和git config配置是完善的  
  还是需要手动去解决conflict
  
### 未完待续...