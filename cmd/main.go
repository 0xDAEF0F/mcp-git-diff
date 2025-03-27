package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	gitops "github.com/0xDAEF0F/mcp-git-diff/internal"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	repoURL = "https://github.com/0xDAEF0F/whistle.git"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Git Diff Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Add git diff tool
	gitDiffTool := mcp.NewTool("git-diff",
		mcp.WithDescription("Perform git diff operations"),
		mcp.WithString("repository-url",
			mcp.Required(),
			mcp.Description("The URL of the repository to diff"),
		),
		mcp.WithNumber("num-days",
			mcp.Required(),
			mcp.Description("The number of days to diff against: 1.0 = 1 day, 2.0 = 2 days, etc."),
		),
	)

	// Add the git diff handler
	s.AddTool(gitDiffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoUrl, ok1 := request.Params.Arguments["repository-url"].(string)
		numDays, ok2 := request.Params.Arguments["num-days"].(float64)

		if !ok1 || !ok2 {
			return nil, errors.New("Invalid arguments")
		}

		time := time.Now().AddDate(0, 0, int(-numDays))
		commits, cleanup, err := gitops.GetRepoCommitsFrom(repoUrl, time)
		defer cleanup()

		if err != nil {
			return nil, errors.New("Failed to retrieve commits")
		}

		lastCommit := commits[0]
		firstCommit := commits[len(commits)-1]

		patch, err := firstCommit.Patch(lastCommit)
		if err != nil {
			return nil, errors.New("Failed to retrieve patch")
		}

		return mcp.NewToolResultText(patch.String()), nil
	})

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
