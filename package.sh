echo "$DOCKER_PASSWORD" | docker login docker.pkg.github.com -u "$DOCKER_USERNAME" --password-stdin 
docker tag stella docker.pkg.github.com/adityaxdiwakar/stella:$TRAVIS_TAG
docker push docker.pkg.github.com/adityaxdiwakar/stella:$TRAVIS_TAG
