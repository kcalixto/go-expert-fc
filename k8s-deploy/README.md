docker build -t kcalixto/go-expert-k8s:latest -f Dockerfile.prod .

docker run --rm -p 8080:8080 kcalixto/go-expert-k8s:latest

kind create cluster --name go-expert-k8s

<!-- careful with this babe -->
kubectl cluster-info --context kind-go-expert-k8s

kubectl apply -f k8s/deployment.yml
kubectl get pods