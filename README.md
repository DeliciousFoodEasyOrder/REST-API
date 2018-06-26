# REST-API Deployment 部署

DFEO项目后台REST API服务部署文档。本服务镜像已打包上传至dockerhub，直接运行`docker-compose.yml`即可部署。

## 使用

```bash
$ sh up.sh
```

产生两个容器，服务器`rest-api`和数据库`db`。
