---
title: Setting up EKS
date: 2019-08-20
for: Myself
---
# EKS learning
## Steps to set up AWS EKS , all in one 

1. run [aws-eks-vpc.yml](https://gist.github.com/aniketchopade/29c1f1bc5b4856dc4aa5a4a9e929885d#file-aws-eks-vpc-yml) in CF 
   this spins up VPC

2. Spin up EKS cluster from AWS Console.
   Note down the kubernetes version (1.12 default as of 06/29/2019)

In any EC2 or on localhost setup:

1. kubectl
    * curl -o kubectl https://amazon-eks.s3-us-west-2.amazonaws.com/1.13.7/2019-06-11/bin/linux/amd64/kubectl
    * chmod +x ./kubectl
    * mkdir -p $HOME/bin && cp ./kubectl $HOME/bin/kubectl && export PATH=$HOME/bin:$PATH
    * echo 'export PATH=$HOME/bin:$PATH' >> ~/.bashrc
    * kubectl version --short --client

2. iam authenticator
    * curl -o aws-iam-authenticator https://amazon-eks.s3-us-west-2.amazonaws.com/1.13.7/2019-06-11/bin/linux/amd64/aws-iam-authenticator
    * chmod +x ./aws-iam-authenticator
    * mkdir -p $HOME/bin && cp ./aws-iam-authenticator $HOME/bin/aws-iam-authenticator && export PATH=$HOME/bin:$PATH
    * echo 'export PATH=$HOME/bin:$PATH' >> ~/.bashrc
    * aws-iam-authenticator help

3. configure kubectl for EKS
    * you need to have updated pip for this
    * curl -O https://bootstrap.pypa.io/get-pip.py
    * python get-pip.py --user
    * pip install awscli --upgrade --user
    * export PATH=$HOME/.local.bin:$PATH
    * echo 'export PATH=$HOME/.local.bin:$PATH' >> ~/.bashrc
    * aws eks update-kubeconfig --name $YOUR_EKS_CLUSTER_NAME
    * kubectl get svc

4. Setup worker node. refer CF template aws-eks-nodegroup.yml
    * this makes sure that correct role is assigned
    * Please follow guide here for pulling correct AMI.
    * Correct AMI (matching with region and kubernetes version is imp. Worker nodes will not be in ready state , if correct AMI not chosen)
    * https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html

5. Configure worker nodes
    * configuration map
    * curl -o aws-auth-cm.yaml https://amazon-eks.s3-us-west-2.amazonaws.com/cloudformation/2019-02-11/aws-auth-cm.yaml
    * this is of kind: ConfigMap

6. Run 
    * Get EC2 worker node's Instance role's ARN 
    * kubectl apply -f aws-auth-cm.yaml

## Dashboard setup 
7. Deploy dashboards
    - Steps (https://docs.aws.amazon.com/eks/latest/userguide/dashboard-tutorial.html)
    - kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/master/aio/deploy/recommended/kubernetes-dashboard.yaml
    - kubectl apply -f https://raw.githubusercontent.com/kubernetes/heapster/master/deploy/kube-    config/influxdb/heapster.yaml
    - kubectl apply -f https://raw.githubusercontent.com/kubernetes/heapster/master/deploy/kube-config/influxdb/influxdb.yaml
    - kubectl apply -f https://raw.githubusercontent.com/kubernetes/heapster/master/deploy/kube-config/rbac/heapster-rbac.yaml
    - create file eks-admin-service-account.yaml
    - kubectl apply -f eks-admin-service-account.yaml
    - kubectl -n kube-system get secret (find out <secret name> for eks-admin user)
    - kubectl describe secret <secretname>
    - get tokn
    - kubectl proxy
    - do tunneling 

8. Open dashboard
   http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login

## kubectl & kubeconfig file & context
 
9. Kubectl & contexts
    - kubectl interacts with kubernetes clusters inconjuction with context(s) defined in ~/.kube/config
    - kubectl config current-context -- this will show current context
    - kubectl config set current-context  docker-for-desktop -- this will set current context
    - override ~/.kube/config with custom kubeconfig file  
        - https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/
        - kubectl --kubeconfig /custom/path/kube.config get pods

10. Kubeconfig file structure
    - clusters
    - contexts
    - current-context
    - users
        - either with client cert & key (for k8 API server to do client auth)
        - or with aws-iam-authenticator.exe

## Kubernetes API services
    - K8 uses declarative API (yml driven)
        - we create an object using CLI or REST to represent what we want the system to do
    -  To get all API resources supported by K8 installation
        - kubectl api-resource -o wide
            - SHORTNAME listed here is used in * kubectl command kubectl <SHORTNAME> * in yml files under "Kind"
            - APIGROUP is used in yml files version
    -  Explain what API requires
        -  kubectl explain <api-name>
        -  kubectl explain ClusterRoleBinding
    -  Taken from https://akomljen.com/kubernetes-api-resources-which-group-and-version-to-use/
