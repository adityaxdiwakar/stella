echo "$DOCKER_PASSWORD" | docker login docker.pkg.github.com -u "$DOCKER_USERNAME" --password-stdin 
docker push adityaxdiwakar/stella
