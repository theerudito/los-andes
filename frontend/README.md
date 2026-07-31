docker build -t istla-frontend .

docker run -d --name react-istla -p 5173:80 istla-frontend
