
from tortoise import fields


class BaseModel():
    created_time = fields.DatetimeField(auto_now_add=True)
    updated_time = fields.DatetimeField(auto_now=True, null=False)