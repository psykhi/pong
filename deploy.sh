make prod
heroku container:login
heroku container:push web -a pong-wasm
heroku container:release web -a pong-wasm

firebase deploy