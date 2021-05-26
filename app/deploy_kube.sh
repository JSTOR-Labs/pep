#!/usr/bin/env bash

#Source this file to get these environment variables.
export SGK_ENVIRONMENT="${1:-test}"
export CLUSTER_NAME="${2:-labs-${SGK_ENVIRONMENT}}"
export KUBECONFIG="$HOME/.kube/kubeconfig.${CLUSTER_NAME}.yaml"
export ASDF_KUBECTL_VERSION=1.13.10
# export HEALTHCHECK_URI="/plants/"
export APP_NAME="pep"
NAMESPACE=pep

case $SGK_ENVIRONMENT in
    test|prod) echo SGK_ENVIRONMENT value $SGK_ENVIRONMENT is valid ;;
    *) echo "SGK_ENVIRONMENT $SGK_ENVIRONMENT not valid" && exit 1;;
esac


AWS_BIN=/usr/local/bin/aws

echo "checking prerequisites:"
prerequisites=( $AWS_BIN kubectl docker curl git aws-iam-authenticator jq )
for i in "${prerequisites[@]}"
do
    if [ ! -x "$(command -v $i)" ]; then
        echo "$i not found. Try:"
        if [[ $i == *"docker"* ]]; then
          echo "brew cask install docker"
        elif [[ $i == *"bin"* ]]; then
          echo "brew install awscli"
        else
          echo "brew install $i"
        fi
        exit 1
    fi
done

if (! docker stats --no-stream 2>/dev/null ); then
    # On Mac OS this would be the terminal command to launch Docker
    echo Docker is not running so starting
    open /Applications/Docker.app
    # Wait until Docker daemon is running and has completed initialisation
    while (! docker stats --no-stream 2>/dev/null ); do
      # Docker takes a few seconds to initialize
      echo "Waiting for Docker to launch..."
      sleep 1
    done
fi

if [ ! -f $KUBECONFIG ]; then

    printf "attempting to update your $KUBECONFIG context for the cluster ${CLUSTER_NAME}...\n"
    mkdir -p "${HOME}/.kube"
    EP=$(${AWS_BIN} eks describe-cluster --name ${CLUSTER_NAME}  --query cluster.endpoint --output text)
    CC=$(${AWS_BIN} eks describe-cluster --name ${CLUSTER_NAME}  --query cluster.certificateAuthority.data --output text)

cat <<EOK > $KUBECONFIG
apiVersion: v1
clusters:
- cluster:
    server: ${EP}
    certificate-authority-data: ${CC}
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: aws
  name: aws
current-context: aws
kind: Config
preferences: {}
users:
- name: aws
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: aws-iam-authenticator
      args:
        - "token"
        - "-i"
        - "${CLUSTER_NAME}"
EOK
    printf "\n...done\n\n"

fi
echo "Setting current kubeconfig $KUBECONFIG context for the cluster ${CLUSTER_NAME} context to $NAMESPACE namespace"
kubectl config set-context aws --namespace=$NAMESPACE

export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-us-east-1}"
GIT_REVISION="$(git rev-parse HEAD)"
IMAGE_NAME="${APP_NAME}"
REPOSITORY_NAMESPACE="labs"
REGISTRY="${AWS_ACCOUNT_ID:-594813696195}.dkr.ecr.${AWS_DEFAULT_REGION:-us-east-1}.amazonaws.com"
REPOSITORY="${REGISTRY}/${REPOSITORY_NAMESPACE}/${IMAGE_NAME}"

if output="$(git status --porcelain)" && [ -z "$output" ]; then
  echo "Git working directory is clean. Checking if image tag already exists"
  # Working directory clean
  set +e
  $AWS_BIN ecr create-repository --repository-name "${REPOSITORY_NAMESPACE}/${IMAGE_NAME}" > /dev/null 2>&1 || true
  IMAGE_META="$( $AWS_BIN ecr describe-images --output json --repository-name=${REPOSITORY_NAMESPACE}/${IMAGE_NAME} --image-ids=imageTag=${GIT_REVISION} 2> /dev/null )"
    if [[ $? == 0 ]]; then
        IMAGE_TAGS="$( echo ${IMAGE_META} | jq '.imageDetails[0].imageTags[0]' -r )"
        echo "${REPOSITORY_NAMESPACE}/${IMAGE_NAME}:${GIT_REVISION} found, so skipping rebuilding and push"
    else
        echo "Building docker image"
        eval "$(/usr/local/bin/aws ecr get-login --region $AWS_DEFAULT_REGION --no-include-email)"


        DIRECTORY="."

        docker build -t "${REPOSITORY}" -t "${REPOSITORY}:${GIT_REVISION}" "${DIRECTORY}"
        echo pushing
        docker push "${REPOSITORY}"

    fi
else
  echo "There are Uncommitted changes. Please commit and try again"
  exit 1
fi


echo "Deploying to $SGK_ENVIRONMENT using image ${REPOSITORY}:${GIT_REVISION}"

cat <<EOF | kubectl -n "${NAMESPACE}" apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${APP_NAME}
  namespace: ${NAMESPACE}
  labels:
    app: ${APP_NAME}
    git: "${GIT_REVISION}"
spec:
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels: # tell deployment which pod to update, should match pod template labels
      app: ${APP_NAME}
  template: # create pods using pod definition in this template
    metadata:
      labels:
        app: ${APP_NAME}
        git: "${GIT_REVISION}"
    spec:
      containers:
      - name: ${APP_NAME}
        image: ${REPOSITORY}:${GIT_REVISION}
        tty: true
        ports:
        - containerPort: 8080
        env:
        - name: THIS_INSTANCE_IP
          valueFrom:
            fieldRef:
              fieldPath: status.hostIP
        - name: THIS_POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: SGK_APP
          value: "${APP_NAME}"
        - name: SGK_ENVIRONMENT
          valueFrom:
            configMapKeyRef:
              name: seq-base
              key: sequoia.k8s.environment
        - name: DEPLOY_ENV
          value: "JSTOR"
        readinessProbe:
          httpGet:
            path: "${HEALTHCHECK_URI}"
            port: 8080
        livenessProbe:
          httpGet:
            path: "${HEALTHCHECK_URI}"
            port: 8080
        resources:
          requests:
            cpu: 100m
            memory: 300Mi
          limits:
            memory: 600Mi
      - name: go-eureka
        image: 594813696195.dkr.ecr.us-east-1.amazonaws.com/playground/go-eureka-eks:cec241f0e6b3a6d2d2152a7f2362de4812180d46
        imagePullPolicy: Always
        env:
        - name: APP_CONT_ID
          value: "K8S"
        - name: APP_NAME
          value: "${APP_NAME}"
        - name: SERVICE_NAME
          value: "${APP_NAME}"
        - name: ENVIRONMENT
          valueFrom:
            configMapKeyRef:
              name: seq-base
              key: sequoia.environment
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: PORT
          value: "8080"
        - name: POD_IP_ADDRESS
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        - name: EUREKA_START_STATUS
          value: "STARTING"
        - name: HEALTHCHECK_URI
          value: "${HEALTHCHECK_URI}"
        resources:
          requests:
            cpu: 100m
            memory: 8Mi
          limits:
            memory: 20Mi
EOF
