# instabug-task
Golang Chat App

This is a Golang chat application powered by a good stack that provides a great tech functionality for a simple chat application might be extended as a full notification system!

---

## Tech Stack
- Docker Containers.
- Redis DB (for cahcing and Streams).
- Logstach (to create Pipeline Synchronization between different data source in this case between Mysql DB and ElasticSearch).
- ElasticSearch (for partial searching for a message accross all chats).
- Mysql (as our main datastore) 
- Kibana (for monitoring elasticsearch queries and visualizing results from elastichsearch)
- Go Application.


<!-- GETTING STARTED -->
## Getting Started

If you want to run this project on locale, take a look to the following steps:

### Prerequisites and Installation

First you should have `docker` and `docker-compose` installed on your machine in order to build our containers

Then you should clone the project using this line of code 
  ```sh
 git clone https://github.com/omarshindy/instabug-task.git
  ```

Then navigate to instabug-task directory from terminal and write the following command
* docker
  ```sh
  docker-compose up --build
  ```
  Note that the container might take up to 3 minutes in the first build due to the healthcheck done in the docker-compose.yaml file to ensure that installation was done in a right way.

# RESTAPIDocs Examples

These examples were taken from project.


## Application Endpoints

* [Create New App](readme/CreateApplication.md) : `POST /api/application/`

* [Update Existing App](readme/UpdateApplications.md) : `POST /api/application/{application_token}`

* [Get All Existing Apps](readme/GetAllApps.md) : `GET /api/applications`


## Chat Endpoint

* [Create New Chat](readme/CreateChat.md) : `POST /api/application/{application_token}/chat`


## Messages Endpoints

* [Create New Message](readme/CreateMessage.md) : `POST /api/application/{application_token}/{chat_number}/message`

* [Get All Existing Messages](readme/GetAllMessages.md) : `GET /api/application/{application_token}/{chat_number}/message`


## Search Endpoint

* [Search Messages](readme/Search.md) : `POST /api/application/{application_token}/chat`
