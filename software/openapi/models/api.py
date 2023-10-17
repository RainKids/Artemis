from tortoise import fields
from tortoise.models import Model

from .basemodel import BaseModel
class Api(Model, BaseModel):
    id = fields.IntField(pk=True)
    name = fields.TextField(null=False)
    content = fields.TextField(null=False)

    class Meta:
        verbose_name = "api"
        verbose_name_plural = verbose_name
        ordering = ["created_time"]