from fastapi import APIRouter,HTTPException
from starlette import status  as HTTPStatus
from common.requestx import request as saleor
from config.config import config
from starlette.requests import Request
from i18n import gettext as _
router = APIRouter()


@router.get("/")
async def query_delivery(request: Request):
    ok, res = await saleor("GET",config.SERVERURI+"/produtcs/delivery",request.headers)
    if ok:
        return res
    raise HTTPException(status_code=HTTPStatus.HTTP_500_INTERNAL_SERVER_ERROR, detail=_("failed"))
