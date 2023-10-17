import os

if os.getenv("DEPLOY_ENV", "") == "dev" or os.getenv("DEPLOY_ENV", "") == "test":
    TORTOISE_ORM = {
        "connections": {
            # Dict format for connection
            "default": "postgres://postgres:123456@127.0.0.1:5432/fastapi-test"
        },
        "apps": {
            "models": {
                # 设置key值“default”的数据库连接
                "default_connection": "default",
                "models": [
                    "aerich.models",
                    "models.api",
                ],
            }
        },
        # 'use_tz': False,
        # 'timezone': 'Asia/Shanghai'
    }
elif os.getenv("DEPLOY_ENV", "") == "prod":
    TORTOISE_ORM = {
        "connections": {
            # Dict format for connection
            "default": f"postgres://postgres:123456@127.0.0.1:5432/fastapi-test"
        },
        "apps": {
            "models": {
                # 设置key值“default”的数据库连接
                "default_connection": "default",
                "models": [
                    "aerich.models",
                    "app.models.quote",
                ],
            }
        },
        # 'use_tz': False,
        # 'timezone': 'Asia/Shanghai'
    }
else:
    TORTOISE_ORM = {
        "connections": {
            # Dict format for connection
            "default": "postgres://postgres:123456@127.0.0.1:5432/fastapi-test"
        },
        "apps": {
            "models": {
                # 设置key值“default”的数据库连接
                "default_connection": "default",
                "models": [
                    "aerich.models",
                    "models.api",
                ],
            }
        },
        # 'use_tz': False,
        # 'timezone': 'Asia/Shanghai'
    }

import logging
import sys

fmt = logging.Formatter(
    fmt="%(asctime)s - %(name)s:%(lineno)d - %(levelname)s - %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
sh = logging.StreamHandler(sys.stdout)
sh.setLevel(logging.DEBUG)
sh.setFormatter(fmt)

# will print debug sql
logger_db_client = logging.getLogger("tortoise.db_client")
logger_db_client.setLevel(logging.DEBUG)
logger_db_client.addHandler(sh)