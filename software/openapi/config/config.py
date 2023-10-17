import os

from pydantic import BaseSettings

from config.db import TORTOISE_ORM


class Settings(BaseSettings):
    REDIS_HOST: str = ""
    REDIS_PORT: str = "6379"
    REDIS_DB: str = "0"

    SECRET_PROJECT_ID: str = ""
    SECRET_ID: str = ""
    SECRET_VERSION_ID: str = ""

    REDIS_EXPIRE: int = 3600 * 24 * 30

    REDIS_TOKEN_EXPIRE: int = 60 * 2
    # 数据库配置
    DATABASE_CONFIG: dict = TORTOISE_ORM
    PAGE_SIZE: int = 20
    MAX_PAGE_SIZE: int = 200

    SERVERURI:str = ""

    class Config:
        env = os.getenv("DEPLOY_ENV", "local")
        env_file = f"{env}.env"

    def load_secrete_manager(self):
        env = os.getenv("DEPLOY_ENV", "local")
        redis_host = os.getenv("REDIS_HOST")
        self.REDIS_HOST = redis_host
        print(
            "配置：",
            self.REDIS_HOST,
        )
        if env in ["dev", "prod"]:
            self.SERVERURI = "http://127.0.0.1:8080/api/v1"
        
        self.SERVERURI = "http://127.0.0.1:8080/api/v1"

config = Settings()
config.load_secrete_manager()
