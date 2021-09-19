# go-container

Getting Filesystem:
```
sudo mkdir /rootfs
cd /rootfs
docker pull ubuntu:18.04
docker save -o ubuntu.tar ubuntu:18.04
ls -la ubuntu.tar
sudo tar -xvf ubuntu.tar
sudo tar -xvf */layer.tar
```


Container from scratch with golang
```
# go run main.go run /bin/bash
# cat /etc/*release
# ps aux 
```

Steps from https://medium.com/adg-vit/creating-your-own-docker-with-go-9a9e978c3918
