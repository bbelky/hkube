package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// import "github.com/mozillazg/request"

// import "golang.org/x/sys"

func main() {

	arg := os.Args[1]

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if arg == "config" {
		fmt.Println("Start loading hkube configuration...")

		// Open our jsonFile
		jsonFile, err := os.Open("config.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}

		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)

		hcloud_token := result["hcloud_token"].(string)
		hcloud_user_sshkey_name := result["hcloud_user_sshkey_name"].(string)
		hcloud_server_type := result["hcloud_server_type"].(string)
		hcloud_location := result["hcloud_location"].(string)
		hcloud_name := result["hcloud_name"].(string)
		hcloud_count := result["hcloud_count"].(string)

		dataResponse := fmt.Sprint(`hcloud_token = "`, hcloud_token, `"`, "\n"+
			`hcloud_user_sshkey_name = "`, hcloud_user_sshkey_name, `"`, "\n"+
			`hcloud_server_type = "`, hcloud_server_type, `"`, "\n"+
			`hcloud_location = "`, hcloud_location, `"`, "\n"+
			`hcloud_name = "`, hcloud_name, `"`, "\n"+
			`hcloud_count = "`, hcloud_count, `"`)

		dataBytes := []byte(dataResponse)
		ioutil.WriteFile("terraform.tfvars", dataBytes, 0644)

		cmd := exec.Command("terraform", "init")
		if err := cmd.Run(); err != nil {
			fmt.Println("Error1:", err)
		}

		if _, err := os.Stat("kubespray"); os.IsNotExist(err) {

			fmt.Println("Downloading kubespray...")
			cmd := exec.Command("git", "clone", "--progress", "https://github.com/kubernetes-sigs/kubespray")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error2:", err)
			}
			fmt.Printf("%s\n", stdout.String())
			fmt.Printf("%s\n", stderr.String())
			stdout.Reset()
			stderr.Reset()

			fmt.Println("Installing requirements...")
			cmd = exec.Command("sudo", "pip", "install", "-r", "requirements.txt")
			cmd.Dir = "kubespray"
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error4:", err)
			}
			fmt.Printf(stdout.String())
			fmt.Printf(stderr.String())

			// cp -rfp inventory/sample inventory/mycluster
			cmd = exec.Command("cp", "-rfp", "inventory/sample", "inventory/mycluster")
			cmd.Dir = "kubespray"
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error4:", err)
			}
			fmt.Printf(stdout.String())
			fmt.Printf(stderr.String())

		} else {

			fmt.Println("Updating kubespray...")
			cmd := exec.Command("git", "-C", "kubespray", "pull")
			cmd.Stdout = &stdout
			if err := cmd.Run(); err != nil {
				fmt.Println("Error3:", err)
			}
			fmt.Println(stdout.String())
		}

		fmt.Println("Configuration loaded successfully")

	}

	if arg == "deploy" {

		fmt.Println("Starting deployment...")

		cmd1 := exec.Command("terraform", "apply", "-auto-approve")
		cmd1.Stdout = &stdout
		cmd1.Stderr = &stderr
		if err := cmd1.Run(); err != nil {
			fmt.Println("Error4:\n", err)
		}
		fmt.Println(stdout.String())
		stdout.Reset()
		stderr.Reset()

		cmd2 := exec.Command("terraform", "output", "public_ip4")
		cmd2.Stdout = &stdout
		cmd2.Stderr = &stderr
		if err := cmd2.Run(); err != nil {
			fmt.Println("Error5:\n", err)
		}

		ips := strings.Fields(stdout.String())
		cips := []string{}
		for _, ip := range ips {
			cips = append(cips, strings.Trim(ip, ","))
		}
		iplist := strings.Join(cips, " ")
		fmt.Println(iplist)

		fmt.Println("Deployment successful!")

		//ansible(iplist)
		ansible(iplist)

	}

	if arg == "destroy" {
		fmt.Println("Destroying kubernetes cluster...")

		cmd := exec.Command("terraform", "destroy", "-auto-approve")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}

		cmd = exec.Command("rm", "-rf", "kubespray")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Cluster deleted")
	}

	if arg == "test" {
		fmt.Println("Testing someting...")

		//createConfig("10.10.1.3 10.10.1.4 10.10.1.5")

		fmt.Println("Tested")
	}

}

func ansible(iplist string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	dataResponse := fmt.Sprint("declare ", "-a ", "IPS=(", iplist, ")", "\n"+
		`CONFIG_FILE=inventory/mycluster/hosts.ini `, "python3 ", "contrib/inventory_builder/inventory.py ", "${IPS[@]}", "\n"+
		`ansible-playbook `, "-i ", `inventory/mycluster/hosts.ini `, "--user=root ", "--become ", "--become-user=root ", "cluster.yml")

	dataBytes := []byte(dataResponse)
	ioutil.WriteFile("kubespray/ansible.sh", dataBytes, 0644)

	cmd := exec.Command("/bin/sh", "ansible.sh")
	cmd.Dir = "kubespray"
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error4:", err)
	}
	fmt.Printf(stdout.String())
	fmt.Printf(stderr.String())

}
