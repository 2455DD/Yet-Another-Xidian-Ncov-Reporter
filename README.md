# Yet-Another-Xidian-Ncov-Report
## 声明
本程序仅用于学习交流，请在符合当地防疫政策的情况下使用。

## 为什么要搞这个
1. 使用Github Action进行自动填报，不用专门租一台服务器来做这事
2. 原有的[Xidian-Ncov-Report](https://github.com/Apache553/xidian-ncov-report)已经归档不再接受PR
3. 学习一下Golang的相关开发

## 使用方法
### Github Action方式
1. Fork本仓库
2. 打开`Setting>Secrets>Action`，分别添加三个键
   1. `STUDENT_ID` : 值为你的学号
   2. `PASSWORD` :值为你的一网通登陆密码
   3. `LOCATION` : 值为你所在的校区,具体选择为:
      1. `xian_south`：在西电北校区
      2. `xian_north`：在西电南校区
      3. `guangzhou`：在西电广研院
      4. `others`: 在其他地方
         - 当你选择这种填报时需要在OPTIONAL中填写地理位置全称
   4. `OPTIONAL`: 当你在iii中选择`others`时填取你的当前所在位置的全称
      - 当选择这种方式时请增加一个`AMAP_KEY`键
      - 例子:`陕西省西安市西沣路兴隆段266号`
   5. `AMAP_KEY` :存储[高德开放平台](https://lbs.amap.com/)API的值
3. 打开`Action`选择当前需要填报的方式

### 手动填报
1. `git clone`本仓库
2. `go build`

---
## TODO
- [ ] 增加健康卡自动申报功能

- [ ] 增加自定义信息功能
