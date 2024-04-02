### Enthusiast girl bot

Created to simplify the deployment process using a telegram bot

## Example
![UI example](./girl.png)

## Build image
1. `docker build -t enthusiast_girl:latest .`
2. `docker tag enthusiast_girl:latest {REGISTRY}/enthusiast_girl:0.0.8`
3. `docker push {REGISTRY}/enthusiast_girl:0.0.8` 

Provide registry credentials if required using `docker login`

## Run image
1. Ensure that your user has private key which needed for stash connection establishment
2. ```
   ssh-keyscan -H -p 7999 -t rsa {HOSTNAME} > ~/.ssh/known_hosts
   sed -i "s/.*ssh-rsa/[{HOSTNAME}]:7999,[{HOSTIP}]:7999 ssh-rsa/" ~/.ssh/known_hosts
   ```
3. Run using docker swarm: `docker stack deploy -c docker-compose.yaml girl`