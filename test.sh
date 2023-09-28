docker build -t gcidp-agent .
docker run --mount type=bind,source=/Users/diyorkhaydarov/Projects/apollo/website,target=/build/context --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock --env GITHUB_REF_NAME=test gcidp-agent
