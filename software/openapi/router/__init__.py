from fastapi import APIRouter

from router.index import router as index_router
from router.api import router as sysapi_router
from router.delivery import router as delivery_router
from router.user import router as user_router
from router.captcha import router as captcha_router
api_router = APIRouter(prefix="/api/v1")
api_router.include_router(index_router, prefix="/index")
api_router.include_router(sysapi_router, prefix="/api")
api_router.include_router(delivery_router, prefix="/delivery")
api_router.include_router(user_router, prefix="/user")
api_router.include_router(captcha_router, prefix="/captcha")
from i18n import lazy_gettext as _
from i18n import set_locale
@api_router.get("/")
def index():
    set_locale("zh")
    return _("您好")
