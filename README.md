# F-Discover

F-Discover is a backend

## Installation

Use [Docker](https://www.docker.com/) to run.
Generate a file **.env** from a file **.env.sample** and file **serviceAccountKey.json**

```bash
#For Production
#Start
docker-compose -f docker-compose.yml up

#Shutdown
docker-compose -f docker-compose.yml down


#For Development to see Docs API in http://localhost
#Start
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

#Shutdown
docker-compose -f docker-compose.yml -f docker-compose.dev.yml down
```