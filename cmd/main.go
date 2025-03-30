package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xDAEF0F/mcp-git-diff/internal/gitops"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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
		mcp.WithString("author-email",
			mcp.Description("The email of the author to diff against"),
		),
		mcp.WithString("branch",
			mcp.Description("The branch to diff against"),
		),
	)

	// Add the git diff handler
	s.AddTool(gitDiffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoUrl, ok1 := request.Params.Arguments["repository-url"].(string)
		numDays, ok2 := request.Params.Arguments["num-days"].(float64)

		if !ok1 || !ok2 {
			return nil, errors.New("Invalid required arguments")
		}

		// Handle optional parameters
		var authorEmail, branch string
		if authorVal, exists := request.Params.Arguments["author-email"]; exists && authorVal != nil {
			authorEmail, _ = authorVal.(string)
		}

		if branchVal, exists := request.Params.Arguments["branch"]; exists && branchVal != nil {
			branch, _ = branchVal.(string)
		}

		opts := gitops.CommitOpts{
			Since: time.Now().AddDate(0, 0, int(-numDays)),
		}

		if authorEmail != "" {
			opts.AuthorEmail = &authorEmail
		}

		if branch != "" {
			opts.Branch = &branch
		}

		diff, err := gitops.GetDiffWithOpts(repoUrl, &opts)

		if err != nil {
			return nil, errors.New("Failed to retrieve commits")
		}

		return mcp.NewToolResultText(diff), nil
	})

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
