from typing import Any

import ujson
from fastapi import status
from fastapi.encoders import jsonable_encoder
from fastapi.responses import JSONResponse, Response

from i18n import lazy_gettext as _


def JsonResponse(data: Any = None, code: int = 0, msg: str = "成功") -> Response:
    return JSONResponse(
        status_code=status.HTTP_200_OK,
        content=jsonable_encoder(
            {
                "code": code,
                "data": data,
                "msg": _(msg),
            }
        ),
    )


class JSONAPIResponse(Response):
    media_type = "application/json"

    def render(
        self,
        content: Any,
    ) -> bytes:
        status_code = self.status_code

        if content:
            if type(content) is dict:
                response = {
                    "code": content.get("code", 0),
                    "data": content.get("data", content or None),
                    "msg": _(content.get("msg", "成功")),
                }
            else:
                response = {"code": 0, "data": content, "msg": _("成功")}
            if not str(status_code).startswith("2"):
                response["code"] = status_code * 100
                response["data"] = None
                try:
                    response["msg"] = (
                        content.get("detail") or list(content.values())[0][0]
                    )
                except Exception:
                    response["data"] = content
                    response["msg"] = _("失败")
        else:
            response = {"code": 0, "data": None, "msg": _("成功")}
        return ujson.dumps(
            response,
            ensure_ascii=False,
            indent=0,
        ).encode()
