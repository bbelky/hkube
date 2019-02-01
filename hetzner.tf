variable "hcloud_token" {}
variable "hcloud_user_sshkey_name" {}
variable "hcloud_server_type" {}
variable "hcloud_location" {}
variable "hcloud_name" {}
variable "hcloud_count" {}

# Configure the Hetzner Cloud Provider
provider "hcloud" {
  token = "${var.hcloud_token}"
}

# create a server
resource "hcloud_server" "kube" {
  count = "${var.hcloud_count}"
  image       = "ubuntu-16.04"
  server_type = "${var.hcloud_server_type}"
  location = "${var.hcloud_location}"
  ssh_keys    = ["${var.hcloud_user_sshkey_name}"]
  name = "${var.hcloud_name}-${count.index}"
}

# get IP
output "public_ip4" {
  value = "${hcloud_server.kube.*.ipv4_address}"
}
