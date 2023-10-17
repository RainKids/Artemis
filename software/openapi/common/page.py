from config.config  import config
async def get_page(request):
    try:
        page = int(request.query_params.get("page", 1))
    except ValueError:
        page = 1
    try:
        pagesize = int(request.query_params.get("page_size", config.PAGE_SIZE))
        if pagesize > 200:
            pagesize = config.MAX_PAGE_SIZE
    except ValueError:
        pagesize = config.PAGE_SIZE
    return page, pagesize