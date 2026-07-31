docker build -t istla-backend .

docker run -d  --name api-istla  -p 9000:9000  istla-backend