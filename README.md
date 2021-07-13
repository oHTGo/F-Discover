# F-Discover

F-Discover is a backend

## Configuration
- Read more detailed instructions in the file **.env.sample**.
- Generate a file **.env** from a file **.env.sample** and file **serviceAccountKey.json**.
## How to run?
Use [Docker](https://www.docker.com/) to run.
```bash
#For Production
#Start
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up

#Rebuild and start
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build

#Shutdown
docker-compose -f docker-compose.yml -f docker-compose.prod.yml down


#For Development
#Start
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

#Rebuild and start
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

#Shutdown
docker-compose -f docker-compose.yml -f docker-compose.dev.yml down
```
## Documents
See Docs API in http://localhost:5055 (default)
