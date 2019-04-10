hkube - Deploy Kubernetes Cluster to Hetzner Cloud
============================================
Automated Kubernetes cluster deployment to Hetzner Cloud based on kubespray. Hkube automatically creates Hetzner Cloud instances and deploy Kubernetes cluster based on kubespray. 

Quick Start
-----------
Hkube uses ansible and terraform, so you need to install them first to your machine. For example, if you use macOS you can do it using homebrew. 

### Homebrew Installation

    /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

### Ansible Installation

    homebrew install ansible

### Terraform Installation

    homebrew install terraform

### Hkube Installation

Download hkube from github:

    git clone https://github.com/bbelky/hkube/

Next you need to add your credentials and k8s settings to config.json file. Copy it from the example first:

    mv config.json.example config.json

After that open the config.json file and add you credentials.

### Kubernetes Deployment

Create configuration files and download kubespray:

    ./hkube config
    
Deploy your cluster:

    ./hkube deploy   
    ...
    Kubernetes cluster deployed! 
    Kubernetes IPs:  116.203.63.47 116.203.49.67 116.203.58.163

### Connect to cluster

Download k8s config file from master server. First IP in the list is your master server. 

    scp root@MASTER_SERVER_IP:/root/.kube/config kubeconf
    
Connect to your k8s cluster:

    # To get list of nodes:
    kubectl --kubeconfig="kubeconf" get nodes
    # To get cluster config:
    kubectl --kubeconfig="kubeconf" cluster-info

## Delete cluster

    ./hkube destroy
    
## More info

How to use Kubespray https://github.com/kubernetes-sigs/kubespray

Contribution
------------
Hkube is OpenSource software and you are free to contribute. Just send your commit and create merge request.

