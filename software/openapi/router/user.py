from fastapi import APIRouter,HTTPException
from starlette import status  as HTTPStatus
from common.requestx import request as saleor
from config.config import config
from starlette.requests import Request
from i18n import gettext as _
router = APIRouter()


@router.post("/login")
async def user_login(request: Request):
    ok, res = await saleor("POST",config.SERVERURI+"/user/login",json=await request.json())
    if ok:
        return res
    raise HTTPException(status_code=HTTPStatus.HTTP_500_INTERNAL_SERVER_ERROR, detail=_("failed"))
