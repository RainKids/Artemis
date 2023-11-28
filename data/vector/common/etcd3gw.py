import threading
import time

import etcd3
from config.config import Settings
class EtcdClient:

    etcd_client:etcd3.Etcd3Client = None
    lease:etcd3.Lease = None

    @classmethod
    async def init_etcd(self, config: Settings):
        self.etcd_client = etcd3.client(host=Settings.ETCD_HOST)

    async def register(self,key,value):
        lease = self.etcd_client.lease(Settings.ETCD_WATCH_TIME_OUT)
        self.etcd_client.put(key=key,value=value,lease=lease)
        t = threading.Thread(target=self.refresh)
        t.start()

    async def refresh(self):
        self.lease.refresh()

    async def unRegister(self,key):
        self.etcd_client.delete(key)

    async def discover(self,key):
        return self.etcd_client.get(key)