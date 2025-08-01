package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/DinethDilhara/glo/internal/formatter"
	"github.com/DinethDilhara/glo/internal/gitexec"
	"github.com/DinethDilhara/glo/internal/models"
	"github.com/spf13/cobra"
)

type BranchConfig struct {
	Format    string
	Tree      bool
	WithDates bool
	Remote    bool
	All       bool
	Graph     bool
}

type BranchService struct {
	gitExec   *gitexec.GitExecutor
	formatter *BranchFormatter
}

type BranchFormatter struct{}

func NewBranchService() *BranchService {
	return &BranchService{
		gitExec:   gitexec.NewGitExecutor(),
		formatter: &BranchFormatter{},
	}
}

func (s *BranchService) FetchBranches(config *BranchConfig) ([]models.Branch, error) {
	return s.gitExec.GetBranches(config.All, config.Remote)
}

func (f *BranchFormatter) FormatOutput(branches []models.Branch, config *BranchConfig) error {
	if len(branches) == 0 {
		fmt.Println("No branches found.")
		return nil
	}

	switch strings.ToLower(config.Format) {
	case "json":
		return f.formatJSON(branches, config)
	case "table":
		return f.formatTable(branches, config)
	case "tree":
		return f.formatTree(branches, config)
	case "graph":
		return f.formatGraph(branches, config)
	case "color", "":
		if config.Tree {
			return f.formatTree(branches, config)
		} else if config.Graph {
			return f.formatGraph(branches, config)
		} else {
			return f.formatColor(branches, config)
		}
	default:
		return fmt.Errorf("unknown format '%s'. Use: color, json, table, tree, or graph", config.Format)
	}
}

func (f *BranchFormatter) formatJSON(branches []models.Branch, config *BranchConfig) error {
	jsonFormatter := formatter.NewJSONFormatter(true)
	fmt.Println(jsonFormatter.FormatBranches(branches))
	return nil
}

func (f *BranchFormatter) formatTable(branches []models.Branch, config *BranchConfig) error {
	colorFormatter := formatter.NewColorFormatter()
	
	fmt.Println(colorFormatter.FormatHeader("Git Branches"))
	fmt.Println()
	
	fmt.Printf("%-20s %-10s %-15s %-50s %s\n", "Branch", "Type", "Last Commit", "Message", "Author")
	fmt.Println(strings.Repeat("-", 120))
	
	for _, branch := range branches {
		typeColor := formatter.ColorGreen
		if branch.IsRemote {
			typeColor = formatter.ColorRed
		}
		if branch.IsCurrent {
			typeColor = formatter.ColorYellow
		}
		
		branchType := "local"
		if branch.IsRemote {
			branchType = "remote"
		}
		if branch.IsCurrent {
			branchType = "current"
		}
		
		fmt.Printf("%-20s %s%-10s%s %-15s %-50s %s\n",
			branch.Name,
			typeColor, branchType, formatter.ColorReset,
			branch.LastCommitDate,
			truncateString(branch.LastCommitMessage, 48),
			branch.LastCommitAuthor)
	}
	
	return nil
}

func (f *BranchFormatter) formatTree(branches []models.Branch, config *BranchConfig) error {
	colorFormatter := formatter.NewColorFormatter()
	
	fmt.Println(colorFormatter.FormatHeader("Git Branch Tree"))
	fmt.Println()
	
	local := make([]models.Branch, 0)
	remote := make([]models.Branch, 0)
	current := ""
	
	for _, branch := range branches {
		if branch.IsCurrent {
			current = branch.Name
		}
		if branch.IsRemote {
			remote = append(remote, branch)
		} else {
			local = append(local, branch)
		}
	}
	
	fmt.Printf("%s📁 Repository%s\n", formatter.ColorBold, formatter.ColorReset)
	fmt.Printf("│\n")
	
	if len(local) > 0 {
		fmt.Printf("├── %sLocal Branches%s\n", formatter.ColorGreen, formatter.ColorReset)
		for i, branch := range local {
			isLast := i == len(local)-1 && len(remote) == 0
			prefix := "│   ├── "
			if isLast {
				prefix = "│   └── "
			}
			
			branchColor := formatter.ColorGreen
			indicator := "  "
			if branch.IsCurrent {
				branchColor = formatter.ColorYellow
				indicator = "* "
			}
			
			fmt.Printf("%s%s%s%s%s", prefix, indicator, branchColor, branch.Name, formatter.ColorReset)
			if config.WithDates {
				fmt.Printf(" %s(%s)%s", formatter.ColorCyan, branch.LastCommitDate, formatter.ColorReset)
			}
			fmt.Println()
		}
	}
	
	if len(remote) > 0 {
		if len(local) > 0 {
			fmt.Printf("│\n")
		}
		fmt.Printf("└── %sRemote Branches%s\n", formatter.ColorRed, formatter.ColorReset)
		for i, branch := range remote {
			isLast := i == len(remote)-1
			prefix := "    ├── "
			if isLast {
				prefix = "    └── "
			}
			
			fmt.Printf("%s%s%s%s", prefix, formatter.ColorRed, branch.Name, formatter.ColorReset)
			if config.WithDates {
				fmt.Printf(" %s(%s)%s", formatter.ColorCyan, branch.LastCommitDate, formatter.ColorReset)
			}
			fmt.Println()
		}
	}
	
	fmt.Println()
	fmt.Printf("%sCurrent branch: %s%s%s\n", formatter.ColorBold, formatter.ColorYellow, current, formatter.ColorReset)
	
	return nil
}

func (f *BranchFormatter) formatGraph(branches []models.Branch, config *BranchConfig) error {
	colorFormatter := formatter.NewColorFormatter()
	
	fmt.Println(colorFormatter.FormatHeader("Git Branch Graph"))
	fmt.Println()
	
	graphData, err := f.getBranchGraph(branches)
	if err != nil {
		return err
	}
	
	f.drawBranchGraph(graphData, config)
	
	return nil
}

type BranchGraphNode struct {
	Name       string
	IsCurrent  bool
	IsRemote   bool
	Commits    []string
	MergePoint string
	Level      int
	Color      string
}

func (f *BranchFormatter) getBranchGraph(branches []models.Branch) ([]BranchGraphNode, error) {
	nodes := make([]BranchGraphNode, 0)
	
	for i, branch := range branches {
		color := formatter.ColorGreen
		if branch.IsRemote {
			color = formatter.ColorRed
		}
		if branch.IsCurrent {
			color = formatter.ColorYellow
		}
		
		node := BranchGraphNode{
			Name:      branch.Name,
			IsCurrent: branch.IsCurrent,
			IsRemote:  branch.IsRemote,
			Level:     i % 4, 
			Color:     color,
		}
		
		nodes = append(nodes, node)
	}
	
	return nodes, nil
}

func (f *BranchFormatter) drawBranchGraph(nodes []BranchGraphNode, config *BranchConfig) {
	fmt.Println("Commit Graph:")
	fmt.Println()
	
	gitExec := gitexec.NewGitExecutor()
	realCommits, err := gitExec.GetCommitGraph(10) 
	if err != nil {
		fmt.Printf("Error getting commit graph: %v\n", err)
		return
	}
	
	for i, commit := range realCommits {
		color := formatter.ColorYellow 
		if len(nodes) > 0 {
			if strings.Contains(commit.Message, "feature") {
				color = formatter.ColorGreen
			} else if strings.Contains(commit.Message, "hotfix") {
				color = formatter.ColorRed
			} else if strings.Contains(commit.Message, "develop") {
				color = formatter.ColorBlue
			}
		}
		
		if strings.Contains(commit.Message, "merge") || strings.Contains(strings.ToLower(commit.Message), "merge") {
			fmt.Printf("  %s│%s\n", color, formatter.ColorReset)
			fmt.Printf("  %s├─╮%s %s%s%s %s\n", 
				color, formatter.ColorReset,
				color, commit.Hash[:8], formatter.ColorReset,
				commit.Message)
			fmt.Printf("  %s│ │%s\n", color, formatter.ColorReset)
		} else {
			fmt.Printf("  %s●%s %s%s%s %s by %s%s%s\n", 
				color, formatter.ColorReset,
				color, commit.Hash[:8], formatter.ColorReset,
				commit.Message,
				formatter.ColorCyan, commit.Author, formatter.ColorReset)
			if i < len(realCommits)-1 {
				fmt.Printf("  %s│%s\n", color, formatter.ColorReset)
			}
		}
	}
	
	fmt.Println()
	
	fmt.Println("Legend:")
	fmt.Printf("  %s●%s Main/current branch   ", formatter.ColorYellow, formatter.ColorReset)
	fmt.Printf("  %s●%s Feature branches   ", formatter.ColorGreen, formatter.ColorReset)
	fmt.Printf("  %s●%s Hotfix branches\n", formatter.ColorRed, formatter.ColorReset)
	fmt.Printf("  %s├─╮%s Merge commit   ", formatter.ColorCyan, formatter.ColorReset)
	fmt.Printf("  %s│%s Commit line\n", formatter.ColorWhite, formatter.ColorReset)
}

func (f *BranchFormatter) formatColor(branches []models.Branch, config *BranchConfig) error {
	for _, branch := range branches {
		prefix := "  "
		color := formatter.ColorGreen
		
		if branch.IsCurrent {
			prefix = "* "
			color = formatter.ColorYellow
		} else if branch.IsRemote {
			color = formatter.ColorRed
		}
		
		fmt.Printf("%s%s%s%s", prefix, color, branch.Name, formatter.ColorReset)
		
		if config.WithDates {
			fmt.Printf(" %s(%s)%s", formatter.ColorCyan, branch.LastCommitDate, formatter.ColorReset)
		}
		
		fmt.Println()
	}
	
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func (s *BranchService) ExecuteBranchCommand(cmd *cobra.Command, args []string) error {
	if !s.gitExec.IsGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	config := &BranchConfig{}
	config.Format, _ = cmd.Flags().GetString("format")
	config.Tree, _ = cmd.Flags().GetBool("tree")
	config.WithDates, _ = cmd.Flags().GetBool("with-dates")
	config.Remote, _ = cmd.Flags().GetBool("remote")
	config.All, _ = cmd.Flags().GetBool("all")
	config.Graph, _ = cmd.Flags().GetBool("graph")
	
	if config.Format == "" {
		config.Format, _ = cmd.Parent().PersistentFlags().GetString("format")
	}

	branches, err := s.FetchBranches(config)
	if err != nil {
		return fmt.Errorf("error fetching branches: %v", err)
	}

	return s.formatter.FormatOutput(branches, config)
}

func runBranchCommand(cmd *cobra.Command, args []string) {
	branchService := NewBranchService()
	
	if err := branchService.ExecuteBranchCommand(cmd, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Show git branches with enhanced visualization",
	Long: `Display git branches with various formatting and visualization options.

The branch command allows you to:
- View branches in different formats (color, table, tree, graph)
- See branch relationships and merge history
- Include commit dates and author information
- Visualize branch structure with ASCII art

Examples:
  glo branch                          # Show branches with colors
  glo branch --tree                   # Show as tree structure
  glo branch --graph                  # Show ASCII commit graph
  glo branch --format=table           # Show as detailed table
  glo branch --with-dates             # Include last commit dates
  glo branch --all                    # Show all branches (local + remote)
  glo branch --remote                 # Show only remote branches`,
	Run: runBranchCommand,
}

func init() {
	rootCmd.AddCommand(branchCmd)

	branchCmd.Flags().StringP("format", "f", "", "Output format: color, table, tree, graph, json")
	branchCmd.Flags().BoolP("tree", "t", false, "Show branches as tree structure")
	branchCmd.Flags().BoolP("graph", "g", false, "Show ASCII commit graph with branches")
	branchCmd.Flags().BoolP("with-dates", "d", false, "Include last commit dates")
	branchCmd.Flags().BoolP("remote", "r", false, "Show only remote branches")
	branchCmd.Flags().BoolP("all", "a", false, "Show all branches (local and remote)")
}
