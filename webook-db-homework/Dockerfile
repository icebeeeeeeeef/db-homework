# 使用scratch基础镜像（最小）
FROM scratch

# 复制本地构建的二进制文件
COPY webook /app/webook

# 设置工作目录
WORKDIR /app

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["/app/webook"]

