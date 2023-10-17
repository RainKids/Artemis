1. 生成初始化数据配置， db.TORTOISE_ORM是上面配置TORTOISE_ORM的路径
aerich init -t config.db.TORTOISE_ORM
2. 生成后会生成一个aerich.ini文件和一个migrations文件夹
3. 初始化数据库
aerich init-db
4. 修改数据模型后生成迁移文件
aerich migrate  在后面加 --name=xxx, 可以指定文件名
5. 执行迁移
aerich upgrade

6. 回退到上一个版本

aerich downgrade


1. extract file
pybabel extract -F i18n/babel.cfg -k lazy_gettext -o i18n/messages.pot app

2. 初始化翻译目录（如果不存在翻译目录需要处理一下）
pybabel init -i i18n/messages.pot -d i18n/ -l zh_CN

3. 同步翻译
pybabel update -i ai18n/messages.pot -d i18n/

4. 编译翻译文件
修改翻译前置文件，比如中文就是：i18n/zh_CN/LC_MESSAGES/messages.po，完成缺少的翻译内容

编译：pybabel compile -d i18n/