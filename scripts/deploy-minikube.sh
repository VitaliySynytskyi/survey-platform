#!/bin/bash
set -e

# Кольори для виводу
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Перевірка, чи встановлено Minikube
if ! command -v minikube &> /dev/null; then
    echo -e "${RED}Minikube не знайдено. Будь ласка, встановіть Minikube.${NC}"
    exit 1
fi

# Перевірка, чи встановлено kubectl
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}kubectl не знайдено. Будь ласка, встановіть kubectl.${NC}"
    exit 1
fi

# Перевірка, чи встановлено Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker не знайдено. Будь ласка, встановіть Docker.${NC}"
    exit 1
fi

# Запуск Minikube, якщо він не запущений
minikube status &> /dev/null || {
    echo -e "${YELLOW}Minikube не запущено. Запускаємо...${NC}"
    minikube start --cpus=4 --memory=8192 --disk-size=20g
}

# Налаштування Docker для використання Docker-демона Minikube
echo -e "${YELLOW}Налаштовуємо Docker для використання Docker-демона Minikube...${NC}"
eval $(minikube docker-env)

# Директорії для розгортання
BACKEND_DIR="./backend"
FRONTEND_DIR="./frontend"
KUBERNETES_DIR="./kubernetes"

# Збірка Docker-образів
echo -e "${YELLOW}Збираємо Docker-образи...${NC}"

# Масив сервісів для збірки
services=(
    "auth_service"
    "user_service"
    "survey_service"
    "survey_taking_service"
    "response_processor_service"
    "analytics_service"
)

# Збираємо образи для мікросервісів
for service in "${services[@]}"; do
    echo -e "${YELLOW}Збираємо образ для ${service}...${NC}"
    docker build -t "survey-platform/${service}:latest" "${BACKEND_DIR}/services/${service}"
done

# Збираємо образ для API Gateway
echo -e "${YELLOW}Збираємо образ для API Gateway...${NC}"
docker build -t "survey-platform/api-gateway:latest" "${BACKEND_DIR}/api_gateway"

# Збираємо образ для Frontend
echo -e "${YELLOW}Збираємо образ для Frontend...${NC}"
docker build -t "survey-platform/frontend:latest" "${FRONTEND_DIR}"

# Створюємо namespace, якщо він ще не існує
kubectl get namespace survey-platform &> /dev/null || {
    echo -e "${YELLOW}Створюємо namespace survey-platform...${NC}"
    kubectl create namespace survey-platform
}

# Застосовуємо ConfigMap та Secret
echo -e "${YELLOW}Застосовуємо ConfigMap та Secret...${NC}"
kubectl apply -f "${KUBERNETES_DIR}/configmaps/survey-platform-config.yaml" -n survey-platform
kubectl apply -f "${KUBERNETES_DIR}/secrets/survey-platform-secrets.yaml" -n survey-platform

# Застосовуємо StatefulSets для баз даних та брокера повідомлень
echo -e "${YELLOW}Розгортаємо інфраструктурні компоненти...${NC}"
kubectl apply -f "${KUBERNETES_DIR}/statefulsets/postgres.yaml" -n survey-platform
kubectl apply -f "${KUBERNETES_DIR}/statefulsets/mongodb.yaml" -n survey-platform
kubectl apply -f "${KUBERNETES_DIR}/statefulsets/rabbitmq.yaml" -n survey-platform
kubectl apply -f "${KUBERNETES_DIR}/statefulsets/consul.yaml" -n survey-platform

# Очікуємо готовності інфраструктурних компонентів
echo -e "${YELLOW}Очікуємо готовності інфраструктурних компонентів...${NC}"
kubectl wait --for=condition=ready pod -l app=postgres -n survey-platform --timeout=180s
kubectl wait --for=condition=ready pod -l app=mongodb -n survey-platform --timeout=180s
kubectl wait --for=condition=ready pod -l app=rabbitmq -n survey-platform --timeout=180s
kubectl wait --for=condition=ready pod -l app=consul -n survey-platform --timeout=180s

# Застосовуємо Deployments та Services для мікросервісів
echo -e "${YELLOW}Розгортаємо мікросервіси...${NC}"
kubectl apply -f "${KUBERNETES_DIR}/deployments/" -n survey-platform
kubectl apply -f "${KUBERNETES_DIR}/services/" -n survey-platform

# Включаємо Ingress у Minikube, якщо він ще не включений
if ! minikube addons list | grep -q "ingress: enabled"; then
    echo -e "${YELLOW}Включаємо Ingress у Minikube...${NC}"
    minikube addons enable ingress
fi

# Застосовуємо Ingress для API Gateway
echo -e "${YELLOW}Налаштовуємо Ingress для API Gateway...${NC}"
kubectl apply -f "${KUBERNETES_DIR}/ingress/api-gateway-ingress.yaml" -n survey-platform

# Отримуємо IP Minikube для доступу до додатку
MINIKUBE_IP=$(minikube ip)

echo -e "${GREEN}Розгортання завершено успішно!${NC}"
echo -e "${GREEN}Додайте наступний запис до вашого файлу /etc/hosts:${NC}"
echo -e "${YELLOW}${MINIKUBE_IP} survey-platform.example.com${NC}"
echo -e "${GREEN}Після цього ви зможете отримати доступ до платформи за адресою:${NC}"
echo -e "${YELLOW}http://survey-platform.example.com${NC}" 