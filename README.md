hkube - Deploy Kubernetes Cluster to Hetzner Cloud
============================================
Automated Kubernetes cluster deployment to Hetzner Cloud based on kubespray. Hkube automatically creates Hetzner Cloud instances and deploy Kubernetes cluster based on kubespray. 

Quick Start
-----------
Hkube uses ansible and terraform, so you need to install them first to your machine. For example, if you use macOS you can do it using homebrew. 

### Homebrew Installation

    # /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

### Ansible Installation

    # homebrew install ansible

### Terraform Installation

    # homebrew install terraform

### Hkube Installation

Download hkube from github:

    # git clone https://github.com/swiftycloud/hkube/

Next you need to add your credentials and k8s settings to config.json file. Copy it from the example first:

    # mv config.json.example config.json

After that open the config.json file and add you credentials.

### Kubernetes Deployment

Create configuration files and download kubespray:

    # ./hkube config
    
Deploy your cluster:

    # ./hkube deploy    

## Delete cluster

    # ./hkube destroy

