FROM python:3.11-alpine

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo Asia/Shanghai > /etc/timezone

# 使用 HTTPS 协议访问容器云调用证书安装
RUN apk add --no-cache ca-certificates

# 指定工作目录
WORKDIR /app

# 拷贝依赖文件
COPY requirements.txt /app/

# pip 源可按需设置国内镜像源以提高下载速度
# 例如使用腾讯云镜像：
RUN pip config set global.index-url https://mirrors.cloud.tencent.com/pypi/simple

# 安装依赖
RUN pip install --no-cache-dir -r requirements.txt

# 拷贝代码
COPY . /app

# 启动命令，默认监听 PORT 环境变量（默认为 8080）
CMD ["sh", "-c", "exec gunicorn -b 0.0.0.0:${PORT:-8080} app:app --access-logfile - --error-logfile -"]
