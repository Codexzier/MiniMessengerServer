docker build -t rpminimessengerserver .

docker tag rpminimessengerserver:latest codexzier/rpminimessengerserver:0.9

docker push codexzier/rpminimessengerserver:0.9