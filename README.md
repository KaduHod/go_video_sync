# Deploy
# Criar Network
```bash
 docker network create sync \
--driver=bridge
```
# Redis Container
```bash
docker pull redis
docker run -p 6379:6379 --hostname video-redis --name video-redis --network sync -d redis
```
# APP Container
```bash
docker build -t KaduHod/sync .
docker run -p 3003:3003 --network sync --name sync KaduHod/sync -d
```
