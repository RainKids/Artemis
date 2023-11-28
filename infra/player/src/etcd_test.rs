fn main() {
    println!("Hello, world!");
}

// 引入依赖库
use etcd_client:: {Client, KeyValue};
use futures:: {StreamExt, TryStreamExt};
use std::sync::Arc;
use tokio::sync::RwLock;
use tonic:: {transport::Server, Request, Response, Status};

// 定义服务接口
pub mod hello_world {
    tonic::include_proto!("helloworld");
}

use hello_world::greeter_server:: {Greeter, GreeterServer};
use hello_world:: {HelloReply, HelloRequest};

// 定义服务实现
# [derive(Debug, Default)]
pub struct MyGreeter {
    // 用于存储服务的地址
    addr: Arc<RwLock<String>>,
}

# [tonic::async_trait]
impl Greeter for MyGreeter {
    // 实现SayHello方法
    async fn say_hello(
        &self,
        request: Request<HelloRequest>,
    ) -> Result<Response<HelloReply>, Status> {
        println!("Got a request: {:?}", request);

        // 从请求中获取名字
        let name = request.into_inner().name;

        // 构造回复消息
        let reply = hello_world::HelloReply {
            message: format!("Hello {}!", name),
        };

        // 返回响应
        Ok(Response::new(reply))
    }
}

// 定义服务注册函数
async fn register_service(client: &Client, key: &str, value: &str) -> Result<(), Box<dyn std::error::Error>> {
    // 使用租约机制，设置服务的过期时间为5秒
    let lease = client.lease_grant(5, None).await?;
    let lease_id = lease.id();

    // 注册服务的键值对，绑定租约
    client
        .kv_put(key, value, Some(lease_id), None)
        .await?;

    // 每隔一定时间（小于过期时间）续约，保持服务的有效性
    let mut interval = tokio::time::interval(tokio::time::Duration::from_secs(1));
    loop {
        interval.tick().await;
        client.lease_keep_alive(lease_id).await?;
    }
}

// 定义服务发现函数
async fn discover_service(client: &Client, key: &str) -> Result<(), Box<dyn std::error::Error>> {
    // 监听键值变化的事件流
    let mut resp = client.watch(key, None, None, None).await?;
    while let Some(event) = resp.message().await? {
        // 根据事件类型处理不同的逻辑
        match event.event_type() {
            // 如果有服务注册，打印服务的地址
            etcd_client::EventType::Put => {
                if let Some(kv) = event.kv() {
                    println!("Service registered: {:?}", kv.value_str()?);
                }
            }
            // 如果有服务注销，打印服务的键
            etcd_client::EventType::Delete => {
                if let Some(kv) = event.prev_kv() {
                    println!("Service deleted: {:?}", kv.key_str()?);
                }
            }
            _ => {}
        }
    }
    Ok(())
}

// 主函数
#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // 连接etcd服务器
    let client = Client::connect(["localhost:2379"], None).await?;

    // 创建服务实例
    let greeter = MyGreeter::default();

    // 获取服务的地址
    let addr = "[::1]:50051".parse().unwrap();

    // 将服务的地址写入服务实例的字段中
    {
        let mut w = greeter.addr.write().await;
        *w = addr.to_string();
    }

    // 创建服务的键，这里使用服务的名称
    let key = "greeter";

    // 获取服务的值，这里使用服务的地址
    let value = &greeter.addr.read().await;

    // 启动一个任务，执行服务注册函数
    tokio::spawn(async move {
        register_service(&client, key, value).await.unwrap();
    });

    // 启动一个任务，执行服务发现函数
    tokio::spawn(async move {
        discover_service(&client, key).await.unwrap();
    });

    // 启动gRPC服务器，监听服务的地址
    Server::builder()
        .add_service(GreeterServer::new(greeter))
        .serve(addr)
        .await?;

    Ok(())
}
