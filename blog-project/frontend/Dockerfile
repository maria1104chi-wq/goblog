# Vue前端Dockerfile
FROM node:18-alpine AS development

WORKDIR /app/frontend

# 复制package.json
COPY package.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY . .

# 开发模式启动
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host"]

# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app/frontend

COPY package.json ./
RUN npm install

COPY . .
RUN npm run build

# 生产镜像
FROM nginx:alpine AS production

# 复制构建产物
COPY --from=builder /app/frontend/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
