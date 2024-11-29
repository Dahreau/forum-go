# forum-go

## 💻 Docker Setup

You can easily run the project using Docker and Docker Compose. Follow these steps to get the application up and running in a Docker container.

#### Step 1: Creating files for persistent database

This step will make sure, that the database will remain unchanged when containers, images are updated.
First creating the directory where the sqlite database will be contained.
````shell
mkdir -p ~/data/forum-go
````
And then creating a empty databsase file
````shell
touch ~/data/forum-go/db.sqlite
````

#### Step 2: Build and Run the Container

To build and run the application, execute the following command:

```shell
docker-compose up --build -d
```

- `up`: This command creates and starts the Docker containers defined in the `docker-compose.yml` file.
- `--build`: This flag ensures that Docker Compose builds the image before starting the container. It’s useful when changes are made to the code or Dockerfile, and you want to rebuild the image. (optional)
- `-d`: This flag runs the container in detached mode, meaning it runs in the background, freeing up your terminal for other tasks. (optional)

Once the command completes, the application will be running in a Docker container, and you can access it on `localhost:8080` (or another port if specified).

#### Step 3: Stop and Remove the Container

When you are done with the container or need to stop it, you can use the following command to stop and remove the running container:

```shell
docker-compose down
```


This will stop the application and clean up the containers, networks, and volumes created by Docker Compose. It's a good way to ensure that the environment is reset when you're done.
