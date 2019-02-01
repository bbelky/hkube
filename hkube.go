package main

import "fmt"
import "encoding/json"
import "os"
import "os/exec"

// import "github.com/mozillazg/request"
import "io/ioutil"

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

		//fmt.Println(dataResponse)

		dataBytes := []byte(dataResponse)
		ioutil.WriteFile("terraform.tfvars", dataBytes, 0644)
		fmt.Println("Configuration loaded successfully")

	}

	if arg == "deploy" {
		fmt.Println("Starting deployment...")

		cmd1 := exec.Command("terraform", "apply", "-auto-approve")
		cmd1.Stdout = os.Stdout
		if err1 := cmd1.Run(); err1 != nil {
			fmt.Println("error1:", err1)
		}
		//fmt.Println("check1:", cmd1.Stdout)

		cmd2 := exec.Command("terraform", "output", "public_ip4")
		iplist := os.Stdout
		if err2 := cmd2.Run(); err2 != nil {
			fmt.Println("error2:", err2)
		}
		fmt.Println("List of IPs:", iplist)

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
