package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// import "github.com/mozillazg/request"

// import "golang.org/x/sys"

func main() {

	arg := os.Args[1]

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

		//fmt.Println(hcloud_token)
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
			fmt.Println("error:", err)
		}
		fmt.Println("Configuration loaded successfully")

	}

	if arg == "deploy" {
		fmt.Println("Starting deployment...")

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd := exec.Command("terraform", "apply", "-auto-approve")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error:\n", err)
		}

		cmd = exec.Command("terraform", "output", "public_ip4")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error:\n", err)
		}

		fmt.Print("List of IPs:\n", stdout.String())

		fmt.Println("Deployment successful!")
	}

	if arg == "destroy" {
		fmt.Println("Destroying kubernetes cluster...")

		cmd := exec.Command("terraform", "destroy", "-auto-approve")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
		fmt.Println(cmd.Stdout)

		fmt.Println("Cluster deleted")
	}

}
