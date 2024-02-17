/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tfdraw",
	Version: version,
	Short:   "Convert tf graph to a diagram format.",
	Long: `This tool takes the terraform graph output and converts it to a mermaid diagram markdown format.
Example:
terraform show --json | tfdraw > diagram.md`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Check if data is being piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			terraformGraph, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
				return
			}

			terraform := decodeJSON(string(terraformGraph))
			// Parse Terraform graph
			nodes, edges := convertTerraformToMermaid(terraform)
			// Generate Mermaid Markdown
			mermaidMarkdown := generateMermaidMarkdown(nodes, edges)

			// Print the Mermaid Markdown diagram
			fmt.Println(mermaidMarkdown)
		} else {
			fmt.Println("No piped data, execute other logic here.")
			// Implement your command logic when there's no piped data
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tfdraw.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func decodeJSON(jsonData string) map[string]interface{} {
	// Trim Byte Order Mark (BOM) to sanitize and convert to byte array for json decoding
	// https://stackoverflow.com/questions/31398044/got-error-invalid-character-%C3%AF-looking-for-beginning-of-value-from-json-unmar
	byteArrayJsonData := bytes.TrimPrefix([]byte(jsonData), []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}

	// Unmarshal the JSON data into a map
	var data map[string]interface{}
	err := json.Unmarshal(byteArrayJsonData, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return data
}

func convertTerraformToMermaid(terraform map[string]interface{}) ([]string, []string) {
	var nodes []string
	var edges []string

	// Accessing nested fields
	rootModule := terraform["values"].(map[string]interface{})["root_module"].(map[string]interface{})
	resources := rootModule["resources"].([]interface{})

	// Accessing individual resources
	for _, resource := range resources {
		resourceMap := resource.(map[string]interface{})
		address := resourceMap["address"].(string)
		nodes = append(nodes, address)
		dependsOn, exists := resourceMap["depends_on"].([]interface{})
		if exists {
			for _, dependent := range dependsOn {
				edge := address + " --> " + dependent.(string)
				edges = append(edges, edge)
			}
		}
	}

	return nodes, edges
}

func generateMermaidMarkdown(nodes, edges []string) string {
	var mermaidMarkdown strings.Builder

	// Start Mermaid graph
	mermaidMarkdown.WriteString("```mermaid\ngraph TD\n")

	// Add nodes
	for _, node := range nodes {
		mermaidMarkdown.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(node)))
	}

	// Add edges
	for _, edge := range edges {
		mermaidMarkdown.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(edge)))
	}

	mermaidMarkdown.WriteString("```\n")

	return mermaidMarkdown.String()
}
