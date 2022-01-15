## Backend API for Jibe

### Docker Start
```dockerfile
docker run -d \
    --network host \
    --restart on-failure:5 \
    --mount type=bind,source=/etc/jibe/,target=/etc/jibe \
    dontpanicdao/jibe-api:lastest
```
