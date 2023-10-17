import httpx

async def request(method, url, params={}, headers={}, http2=False, **kwargs):
    """
    对 request 做一个简单封装
    """
    try:
        async with httpx.AsyncClient(http2=http2,timeout=120) as client:
            r = None
            if method == "GET":
                r = await client.request(method, url, params=params, headers=headers)
            else:
                r = await client.request(
                    method, url, data=params, headers=headers, **kwargs
                )
            if r.status_code == 204:
                return True, dict()
            elif r.status_code != 200:
                return False, dict(
                    result="error",
                    message=f"response code: {r.status_code}, {str(r.content)}",
                )

            result = r.json()
          
            return True, result

    except httpx.TimeoutException:
        return False, dict(result="error", message="request timeout")
    except httpx.RequestError as e:
        return False, dict(result="error", message=str(e))
    except Exception as e:
        return False, dict(result="error", message=str(e))
