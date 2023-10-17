import json
import logging
import sys

from redis import AuthenticationError, Redis

from config.config import Settings

MatomoTokenExpire = 60 * 60 * 24 * 15


class RedisClinet:
    redis: Redis = None

    def init_redis(self, config: Settings):
        try:
            self.redis = Redis(
                host=config.REDIS_HOST,
                port=config.REDIS_PORT,
                db=config.REDIS_DB,
                decode_responses=True,
                health_check_interval=30,
            )
            # TODO: 这个过期时间需要修改
            self.token_expire = config.REDIS_TOKEN_EXPIRE

            if not self.redis.ping():
                logging.error("连接 redis 超时")
                sys.exit()
        except (AuthenticationError, Exception) as e:
            logging.error("连接 redis 异常 : ", e)
            sys.exit()

    def close(self):
        self.redis.close()


cache_client = RedisClinet()
