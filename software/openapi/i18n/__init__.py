import os

from babel import Locale
from babel.support import Translations
from starlette_context import context

from common.lazy_string import LazyString

DEFAULT_LOCALE = Locale.parse("en")
DOMAIN = "messages"


TRANSLATIONS_DICT = {
    # 获取当前文件所在目录，翻译文件就在当前目录下
    "zh-cn": "zh_CN",
    "zh_CN": "zh_CN",
    "zh_Hans_CN": "zh_CN",
    "en": "en",
}

translations_locale = Translations.load(
    os.path.split(os.path.realpath(__file__))[0], locales=["en"], domain=DOMAIN
)


def _get_current_context():
    return context if context.exists() else None


def get_locale():
    rv = context.data.get(
        "Accept-Languages"
    )  # might add some custom locale selector here, i.e: query from database
    if rv is None:
        locale = DEFAULT_LOCALE
    else:
        locale = Locale.parse(rv)
        # setattr(ctx, "babel_locale", locale)
    return locale


def set_locale(locale: str):
    global translations_locale
    translations_locale = Translations.load(
        os.path.split(os.path.realpath(__file__))[0],
        locales=[TRANSLATIONS_DICT.get(locale), "en"],
        domain=DOMAIN,
    ) or Translations.load(
        os.path.split(os.path.realpath(__file__))[0], locales=["en"], domain=DOMAIN
    )


def gettext(msg: str):
    if not msg:
        return ""
    if TRANSLATIONS_DICT.get(context.data.get("Accept-Languages")):
        translations = Translations()
        catalog = Translations.load(
            os.path.split(os.path.realpath(__file__))[0],
            locales=[TRANSLATIONS_DICT.get(context.data.get("Accept-Languages"), "en")],
            domain=DOMAIN,
        )
        translations.merge(catalog)
        if hasattr(catalog, "plural"):
            translations.plural = catalog.plural
        return translations.ugettext(msg)
    else:
        return translations_locale.ugettext(msg)


def lazy_gettext(msg: str):
    return str(LazyString(gettext, msg))


def t(msg, language):
    if not msg:
        return ""
    if TRANSLATIONS_DICT.get(language):
        translations = Translations()
        catalog = Translations.load(
            os.path.split(os.path.realpath(__file__))[0],
            locales=[TRANSLATIONS_DICT.get(language, "en")],
            domain=DOMAIN,
        )
        translations.merge(catalog)
        if hasattr(catalog, "plural"):
            translations.plural = catalog.plural
        return translations.ugettext(msg)
    else:
        return translations_locale.ugettext(msg)
