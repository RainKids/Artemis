import logging
import traceback
import uvicorn
from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from starlette import status
from starlette.exceptions import HTTPException
from starlette.middleware import Middleware
from starlette.middleware.cors import CORSMiddleware
from starlette.requests import Request
from starlette.responses import Response
from starlette_context.middleware import RawContextMiddleware
from starlette_context.plugins.base import Plugin
from tortoise.contrib.fastapi import register_tortoise

from config.config import config
from common.cache import cache_client
from common.response import JSONAPIResponse
from router import api_router



class AcceptLanguagesHeaderPlugin(Plugin):
    key = "Accept-Languages"


middleware = [
    Middleware(
        RawContextMiddleware,
        plugins=(AcceptLanguagesHeaderPlugin(),),
    )
]


def register_hook(app: FastAPI) -> None:
    """
    请求响应拦截 hook
    https://fastapi.tiangolo.com/tutorial/middleware/
    :param app:
    :return:
    """

    @app.middleware("http")
    async def logger_request(request: Request, call_next) -> Response:
        # https://stackoverflow.com/questions/60098005/fastapi-starlette-get-client-real-ip
        # logger.info(f"访问记录:{request.method} url:{request.url}\nheaders:{request.headers}\nIP:{request.client.host}")
        response = await call_next(request)
        return response


def register_init(app: FastAPI) -> None:
    """
    初始化连接
    :param app:
    :return:
    """

    @app.on_event("startup")
    async def init_connect():
        # 连接数据库
        register_tortoise(
            app,
            config=config.DATABASE_CONFIG,
            # modules=settings.DATABASE_CONFIG.get("apps").get("models"),
            generate_schemas=True,  # True 表示连接数据库的时候同步创建表
            add_exception_handlers=True,
        )
        # logger.info("start server and register_tortoise")
        # 连接redis
        # redis_client.init_redis_connect()

        # 初始化 apscheduler
        # schedule.init_scheduler()

    @app.on_event("shutdown")
    async def shutdown_connect():
        """
        关闭
        :return:
        """
        # schedule.shutdown()
        # logger.info('stop server')


def register_exception(app: FastAPI) -> None:
    @app.exception_handler(HTTPException)
    async def http_exception_handler(request: Request, exc: HTTPException):
        return JSONAPIResponse({"detail": str(exc.detail)}, status_code=exc.status_code)

    @app.exception_handler(RequestValidationError)
    async def validation_exception_handler(
        request: Request, exc: RequestValidationError
    ):
        return JSONAPIResponse(
            status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
            content={"detail": exc.errors(), "body": exc.body},
        )


async def catch_exceptions_middleware(request: Request, call_next):
    try:
        return await call_next(request)
    except Exception:
        # you probably want some kind of logging here
        logging.error(traceback.format_exc())
        return Response("Internal dao error", status_code=500)


def create_fast_app():
    app = FastAPI(
        title=__name__,
        middleware=middleware,
        # default_response_class=JSONAPIResponse,
    )

    app.middleware("http")(catch_exceptions_middleware)

    # origins = ["https://f3c8-103-206-188-68.ngrok.io"]
    origins = ["*"]
    app.add_middleware(
        CORSMiddleware,
        allow_origins=origins,
        allow_credentials=True,
        allow_methods=["*"],  # 设置允许跨域的http方法，比如 get、post、put等。
        allow_headers=["*"],  # 允许跨域的headers，可以用来鉴别来源等作用。
    )

    app.include_router(api_router)
    register_init(app)
    register_exception(app)
    register_hook(app)
    return app


app = create_fast_app()


@app.on_event("startup")
def create_redis():
    cache_client.init_redis(config)


@app.on_event("shutdown")
def close_redis():
    cache_client.close()

