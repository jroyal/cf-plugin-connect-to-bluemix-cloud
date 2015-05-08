package main

import (
	"fmt"
  "strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/cloudfoundry/cli/plugin"
)

type OpenstackSignin struct{}

func (c *OpenstackSignin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "connect-to-bluemix-cloud" {
		fmt.Println("Gathering cf information to determine which cloud to use")
		fmt.Println("========================================================")

    // Get the current Org and Space
    output, err := cliConnection.CliCommand("t")
    if err != nil {
  		fmt.Println("PLUGIN ERROR: Error from CliCommand: ", err)
  	}
    org := strings.TrimSpace(strings.Split(output[3], ":")[1])
    space := strings.TrimSpace(strings.Split(output[4], ":")[1])

		// Get the UUID for the org and space
    output, err = cliConnection.CliCommand("org", org, "--guid")
    if err != nil {
  		fmt.Println("PLUGIN ERROR: Error from CliCommand: ", err)
  	}
    org_uuid := strings.TrimSpace(output[0])
    output, err = cliConnection.CliCommand("space", space, "--guid")
    if err != nil {
  		fmt.Println("PLUGIN ERROR: Error from CliCommand: ", err)
  	}
    space_uuid := strings.TrimSpace(output[0])

		// Get the Auth Token
		output, err = cliConnection.CliCommand("oauth-token")
    if err != nil {
  		fmt.Println("PLUGIN ERROR: Error from CliCommand: ", err)
  	}
		bearer_token := strings.TrimSpace(output[3])

		// Figure out which cloud we are focusing on
		url := "https://clouds.stage1.ng.bluemix.net/v1/clouds/?ace_config={" +
			"\"orgGuid\":\""+ org_uuid+"\"," +
			"\"spaceGuid%22:\""+space_uuid+"\"}"

		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", bearer_token)
		resp, _ := client.Do(req)
		if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
        panic(err)
    }

		// Choose the cloud we want
		fmt.Println("\n\n")
		fmt.Println("Clouds available to this space and org")
		fmt.Println("======================================")
		clouds := data["clouds"].([]interface{})
		for key, value := range clouds {
	    fmt.Println(key+1, value.(map[string]interface{})["cloud_name"])
		}
		var cloud_key int
		fmt.Printf("Choose cloud number: ")
		fmt.Scanf("%s", &cloud_key)
		cloud := data["clouds"].([]interface{})[cloud_key].(map[string]interface{})

		region := cloud["region"]
		auth_url := cloud["auth_url"]
		tenant := cloud["tenant"]
		cloud_id := cloud["cloud_id"]

		// Get the username and password for this cloud for the logged in user
		url = "https://clouds.stage1.ng.bluemix.net/v1/clouds/"+cloud_id.(string)+"/credential"
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", bearer_token)
		resp, _ = client.Do(req)
		if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ = ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &data); err != nil {
        panic(err)
    }
		credentials := data["data"].(map[string]interface{})
		user_pass := credentials["cloudpwd"]
		user_name := credentials["clouduser"]

		fmt.Println("\n\n")
		fmt.Println("Copy and paste the following to allow CLI access to your cloud")
		fmt.Println("==============================================================")
		fmt.Printf("export OS_AUTH_URL=%s;export OS_TENANT_NAME=%s;" +
								"export OS_USERNAME=%s; export OS_PASSWORD='%s';"+
								"export OS_REGION_NAME=%s\n\n", auth_url, tenant, user_name, user_pass, region)


		fmt.Println("Successfully ran the connect-to-bluemix-cloud plugin")
	}
}

func (c *OpenstackSignin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "connect-to-bluemix-cloud",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "connect-to-bluemix-cloud",
				HelpText: "Export variables that allow you to connect to your apps cloud",
				UsageDetails: plugin.Usage{
					Usage: "connect-to-bluemix-cloud - Exports variables for openstack cli clients\n   cf connect-to-bluemix-cloud",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(OpenstackSignin))
}
