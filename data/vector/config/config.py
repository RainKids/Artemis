import os

from pydantic import BaseSettings
class Settings(BaseSettings):
    ETCD_HOST: str = ""
    ETCD_PORT: int = 2379
    ETCD_WATCH_TIME_OUT  = 10
    class Config:
        env = os.getenv("DEPLOY_ENV", "local")
        env_file = f"{env}.env"

    def load_secrete_manager(self):
        etcd_host = os.getenv("ETCD_HOST")
        self.ETCD_HOST = etcd_host
        print(
            "配置：",
            self.ETCD_HOST,
        )
config = Settings()
config.load_secrete_manager()
