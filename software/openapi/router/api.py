from fastapi import APIRouter, Depends
from common.response import JsonResponse
from schemas.api import ApiFilter, ApiPageResult
from models.api import Api
from common.db import use_pagination

router = APIRouter()

@router.get("/",summary="api列表",response_model=ApiPageResult)
async def query_api(query: ApiFilter = Depends()):
    data = await use_pagination(query,Api)
    return JsonResponse(data)

@router.get("/{id}", summary="api")
async def get_api(id:int):
    api_obj = await Api.get_or_none(id=id)
    return JsonResponse(api_obj)
