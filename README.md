# F-Discover

F-Discover is a backend

## Installation

Use [Docker](https://www.docker.com/) to run.
Generate a file **.env** from a file **.env.sample** and file **serviceAccountKey.json**.
Read more detailed instructions in the file **.env.sample**.

See Docs API in http://localhost:5055 (default)

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