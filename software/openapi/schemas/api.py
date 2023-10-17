from schemas.schemas import PageResult,Filter
from models.api import Api
from pydantic import Field
from tortoise.contrib.pydantic import pydantic_model_creator


ApiOut = pydantic_model_creator(Api,name="ApiOut")

ApiPageResult = PageResult[list[ApiOut]]


class ApiFilter(Filter):
    name__icontains: str = Field("", alias="name")