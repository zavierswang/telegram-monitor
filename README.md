# TRC20钱包转帐事件监听

### 功能
* 🔥 支持高并发秒级监听TRC20金额变动
* 🔥 支持自定义标签
* 🔥 支持为特定地址添加头像 🎉
* 🔥 支持群组通知 🎉
* 🔥 完善查询钱包收支统计交互
* 🔥 完善管理员授权等多种操作指令
* 🔥 新用户默认可试用3天

### 更多程序
* [telegram-trx](https://github.com/zavierswang/telegram-trx) **TRX兑换机器人**
* [telegram-monitor](https://github.com/zavierswang/telegram-monitor) **TRC20钱包事件监听机器人**
* [telegram-search](https://github.com/zavierswang/telegram-search) **导航机器人**（可支持全网搜索，API收费有点小贵）
* [telegram-premium](https://github.com/zavierswang/telegram-premium) **Telegram Premium自动充值机器人**
* [telegram-replay](https://github.com/zavierswang/telegram-replay) **双向机器人**
* [telegram-energy](https://github.com/zavierswang/telegram-energy) **TRON能量租凭机器人**
* [telegram-proto](https://github.com/zavierswang/telegram-proto) **Telegram协议号机器人**


### 部署
* 本程序基于`Telegram Bot`，分为主程序`telegram-monitor`和`telegram-scanner`
* 确保机器人部署服务可以访问外网`telegram.org`
* 使用自己的`telegram`生成一个机器人，并获取到`token`
* `telegram-scanner`分布式部署到不同服务器，不同服务器应该对应不同的`tron_scan_key`和`tron_grid_key`
* 部署`RabbitMQ`服务
* 配置文件`telegram-monitor.yaml.example`改名为`telegram-monitor.yaml`, 修改建议配置项
* 配置文件`telegram-scanner.yaml.example`改名为`telegram-scanner.yaml`, 修改建议配置项

### 功能演示
* **添加地址**
  
    ![add_address.png](https://github.com/zavierswang/telegram-monitor/blob/main/img/add_address.png)
  * 输入TRC20地址，如果是新地址，可以添加/修改备注/头像操作
  * 添加完成后，您可以使用 **钱包概览** 查看该地址的近两个月的数据统计信息，主要用于统计一段时间内的收支情况
  * 还可以通过 **钱包管理** 进行快速删除

* **查询详情**

    ![list_address.png](https://github.com/zavierswang/telegram-monitor/blob/main/img/list_address.png) 
    ![overview.png](https://github.com/zavierswang/telegram-monitor/blob/main/img/overview.png)
    ![statistics.png](https://github.com/zavierswang/telegram-monitor/blob/main/img/statistics.png)

* **事件通知**

    ![notice.png](https://github.com/zavierswang/telegram-monitor/blob/main/img/notice.png)



> **注意：**
> * 不支持交易所转帐事件监听
> * 如果出现`网络错误`，请增加`telegram-scanner`服务器
> * 对linux不熟悉的给点打赏手把手教学🤭
> * 配置文件中的`license`配置请找 [🫣我](https://t.me/tg_llama) 拿~


