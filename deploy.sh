#aws ecr get-login-password --region eu-north-1 --profile perso | docker login --username AWS --password-stdin 966537388583.dkr.ecr.eu-north-1.amazonaws.com
#docker build -t pong .
#docker tag pong:latest 966537388583.dkr.ecr.eu-north-1.amazonaws.com/pong:latest
#docker push 966537388583.dkr.ecr.eu-north-1.amazonaws.com/pong:latest
#heroku login

heroku container:login
heroku container:push web -a pong-wasm
heroku container:release web -a pong-wasm