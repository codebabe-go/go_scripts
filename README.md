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
     * go build ${project_path}/gitgo.go
  4. 运行
     * 假设你go build 时候的路径为build_path, 生成的二进制文件叫做 *nix(gitgo), windows(gitgo.exe)
     * Windows: 为build生成的gitgoexe文件配置path环境变量, 然后在你存在git repository的目录下执行 gitgo.exe "comment" "push_branch"
     * *nix: .bash_profile "alias gitgo=${build_path}/gitgo", 直接使用gitgo "comment" "push_branch"
  
## 缺点
* 不去获取本地的user信息, 需要保障的前提是本地git命令和git config配置是完善的  
* 无法自动去解决冲突, 需要手动去解决conflict
  
## 后续
* 程序自动解决一些简单的问题正在开发, 加上skip机制

### 未完待续...