package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/MuktadirHassan/bunny-cli/bunny"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	folderPath       string
	filePath         string
	concurrencyLimit int
	timeout          time.Duration
	failFast         bool
	pullZoneName     string
	fileInputPath    string
	storageAccessKey string
	storageZoneName  string
	storageZoneHost  string
)

var RootCmd = &cobra.Command{
	Use:   "bunny-cli",
	Short: "A CLI tool to interact with Bunny.net",
	Long:  `A CLI tool to perform tasks such as uploading files, purging cache, and more using Bunny.net APIs.`,
}

var uploadFolderCmd = &cobra.Command{
	Use:   "upload-folder",
	Short: "Upload a folder concurrently",
	Long:  `Upload all files in a folder concurrently using a worker pool with configurable concurrency and timeout and retry logic.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uploading folder:", folderPath)
		fmt.Println("Concurrency limit:", concurrencyLimit)
		fmt.Println("Timeout:", timeout)
		fmt.Println("Fail fast:", failFast)

		// Call your UploadFolder function here
		err := bunny.UploadFolder(folderPath, concurrencyLimit, timeout, failFast, storageAccessKey, storageZoneName, storageZoneHost)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Upload completed successfully!")
	},
}

var uploadFileCmd = &cobra.Command{
	Use:   "upload-file",
	Short: "Upload a single file",
	Long:  `Upload a single file to Bunny.net.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uploading file:", filePath)
		fmt.Println("Not implemented yet...")
	},
}

var purgeCacheFullCmd = &cobra.Command{
	Use:   "purge-cache-full",
	Short: "Purge the full pull zone cache",
	Long:  `Purge the full cache for a specified pull zone in Bunny.net.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Purging full cache for pull zone:", pullZoneName)

		// Call your PurgeCacheFull function here
		err := bunny.CdnFullCachePurge(pullZoneName)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Cache purged successfully!")

	},
}

var purgeCacheURLCmd = &cobra.Command{
	Use:   "purge-cache-url",
	Short: "Purge cache for URLs listed in a file",
	Long:  `Purge cache for specific URLs listed in a file for a Bunny.net pull zone.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Purging cache for URLs in file:", fileInputPath)
		fmt.Println("Not implemented yet...")
	},
}

var docsCmd = &cobra.Command{
	Use:   "gen-docs",
	Short: "Generate CLI documentation",
	Long:  `Generate documentation for the CLI tool in Markdown format.`,
	Run: func(cmd *cobra.Command, args []string) {
		outputDir := "./docs" // Specify the directory to save docs
		fmt.Println("Generating documentation...")

		err := doc.GenMarkdownTree(RootCmd, outputDir)
		if err != nil {
			fmt.Println("Error generating documentation:", err)
			os.Exit(1)
		}
		fmt.Println("Documentation generated in", outputDir)
	},
}

func init() {
	// Add flags for upload-folder command
	uploadFolderCmd.Flags().StringVarP(&folderPath, "folder", "f", "", "Path to the folder to upload (required)")
	uploadFolderCmd.Flags().IntVarP(&concurrencyLimit, "concurrency", "c", 10, "Number of concurrent workers")
	uploadFolderCmd.Flags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "Timeout for each file upload")
	uploadFolderCmd.Flags().BoolVarP(&failFast, "fail-fast", "F", true, "Enable fail-fast mode")
	uploadFolderCmd.Flags().StringVarP(&storageAccessKey, "access-key", "a", "", "Storage access key (required)")
	uploadFolderCmd.Flags().StringVarP(&storageZoneName, "zone-name", "z", "", "Storage zone name (required)")
	uploadFolderCmd.Flags().StringVarP(&storageZoneHost, "zone-host", "H", "sg.storage.bunnycdn.com", "Storage zone host")
	uploadFolderCmd.MarkFlagRequired("folder")
	uploadFolderCmd.MarkFlagRequired("access-key")
	uploadFolderCmd.MarkFlagRequired("zone-name")

	// Add flags for upload-file command
	uploadFileCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file to upload (required)")
	uploadFileCmd.MarkFlagRequired("file")

	// Add flags for purge-cache-full command
	purgeCacheFullCmd.Flags().StringVarP(&pullZoneName, "pull-zone", "p", "", "Name of the pull zone to purge (required)")
	purgeCacheFullCmd.MarkFlagRequired("pull-zone")

	// Add flags for purge-cache-url command
	purgeCacheURLCmd.Flags().StringVarP(&fileInputPath, "file", "f", "", "Path to the file containing URLs to purge (required)")
	purgeCacheURLCmd.MarkFlagRequired("file")

	// Add subcommands to the root command
	RootCmd.AddCommand(uploadFolderCmd)
	RootCmd.AddCommand(uploadFileCmd)
	RootCmd.AddCommand(purgeCacheFullCmd)
	RootCmd.AddCommand(purgeCacheURLCmd)

	// Add gen-docs command
	RootCmd.AddCommand(docsCmd)
}
