"""Plugin system for tldr++."""

import os
import subprocess
import tempfile
from abc import ABC, abstractmethod
from pathlib import Path
from typing import List, Optional

from tldrpp.cache import Example, Page


class Plugin(ABC):
    """Base class for tldr++ plugins."""
    
    @abstractmethod
    def name(self) -> str:
        """Return the plugin name."""
        pass
    
    @abstractmethod
    def description(self) -> str:
        """Return the plugin description."""
        pass
    
    @abstractmethod
    def execute(self, args: List[str]) -> None:
        """Execute the plugin with given arguments."""
        pass


class SubmitPlugin(Plugin):
    """Plugin for submitting examples to tldr-pages."""
    
    def __init__(self, page: Page, example: Example) -> None:
        """Initialize submit plugin."""
        self.page = page
        self.example = example
    
    def name(self) -> str:
        """Return plugin name."""
        return "submit"
    
    def description(self) -> str:
        """Return plugin description."""
        return "Submit example to tldr-pages repository"
    
    def execute(self, args: List[str]) -> None:
        """Execute submit plugin."""
        if not args:
            raise ValueError("No command specified")
        
        command = args[0]
        
        if command == "init":
            self._init_submission()
        elif command == "validate":
            self._validate_example()
        elif command == "create-pr":
            self._create_pull_request()
        else:
            raise ValueError(f"Unknown command: {command}")
    
    def _init_submission(self) -> None:
        """Initialize a new submission."""
        print("Initializing tldr-pages submission...")
        print(f"Page: {self.page.name} ({self.page.platform})")
        print(f"Example: {self.example.description}")
        print(f"Command: {self.example.command}")
        print()
        
        # Check if git is available
        if not self._is_git_available():
            raise RuntimeError("git is not available. Please install git to submit to tldr-pages")
        
        # Check if gh CLI is available
        if not self._is_github_cli_available():
            print("Warning: GitHub CLI (gh) is not available.")
            print("You'll need to manually create a pull request.")
        
        # Create submission directory
        submission_dir = Path(tempfile.gettempdir()) / "tldrpp-submission"
        submission_dir.mkdir(parents=True, exist_ok=True)
        
        # Generate markdown content
        content = self._generate_markdown()
        content_file = submission_dir / f"{self.page.name}.md"
        content_file.write_text(content, encoding="utf-8")
        
        print(f"Submission files created in: {submission_dir}")
        print("Next steps:")
        print("1. Review the generated markdown file")
        print("2. Run 'tldrpp plugin submit validate' to check for issues")
        print("3. Run 'tldrpp plugin submit create-pr' to create a pull request")
    
    def _validate_example(self) -> None:
        """Validate the example against tldr-pages standards."""
        print("Validating example against tldr-pages standards...")
        
        issues = []
        
        # Check description length
        if len(self.example.description) > 80:
            issues.append("Description is too long (>80 characters)")
        
        # Check command length
        if len(self.example.command) > 100:
            issues.append("Command is too long (>100 characters)")
        
        # Check for common issues
        if "sudo" in self.example.command:
            issues.append("Avoid using 'sudo' in examples")
        
        if "&&" in self.example.command:
            issues.append("Avoid chaining commands with '&&'")
        
        # Check placeholder usage
        for placeholder in self.example.placeholders:
            if not placeholder.name:
                issues.append("Empty placeholder name found")
            if len(placeholder.name) > 20:
                issues.append(f"Placeholder name '{placeholder.name}' is too long")
        
        if not issues:
            print("✓ Example validation passed!")
            return
        
        print("✗ Validation issues found:")
        for issue in issues:
            print(f"  - {issue}")
        
        raise RuntimeError(f"Validation failed with {len(issues)} issues")
    
    def _create_pull_request(self) -> None:
        """Create a pull request to tldr-pages."""
        print("Creating pull request to tldr-pages...")
        
        # Check if gh CLI is available
        if not self._is_github_cli_available():
            raise RuntimeError("GitHub CLI (gh) is not available. Please install it or create a PR manually")
        
        # Generate branch name
        branch_name = f"tldrpp-{self.page.name}-{self.page.platform}"
        
        # Create markdown content
        content = self._generate_markdown()
        
        # Create a temporary file for the content
        with tempfile.NamedTemporaryFile(mode='w', suffix='.md', delete=False) as temp_file:
            temp_file.write(content)
            temp_file_path = temp_file.name
        
        try:
            # Create PR using gh CLI
            title = f"Add example for {self.page.name} ({self.page.platform})"
            body = (f"This PR adds a new example for the `{self.page.name}` command on the `{self.page.platform}` platform.\n\n"
                   f"Example: {self.example.description}\n\n"
                   f"Command: `{self.example.command}`")
            
            subprocess.run([
                "gh", "pr", "create",
                "--repo", "tldr-pages/tldr",
                "--title", title,
                "--body", body,
                "--file", temp_file_path
            ], check=True)
            
            print("✓ Pull request created successfully!")
        
        finally:
            # Clean up temporary file
            os.unlink(temp_file_path)
    
    def _generate_markdown(self) -> str:
        """Generate markdown content for the submission."""
        content = []
        
        # Title
        content.append(f"# {self.page.name}\n")
        
        # Description
        content.append(f"> {self.page.description}.\n")
        
        # Example
        content.append(f"- {self.example.description}:")
        content.append(f"  `{self.example.command}`")
        
        return "\n".join(content)
    
    def _is_git_available(self) -> bool:
        """Check if git is available."""
        try:
            subprocess.run(["git", "--version"], capture_output=True, check=True)
            return True
        except (subprocess.CalledProcessError, FileNotFoundError):
            return False
    
    def _is_github_cli_available(self) -> bool:
        """Check if GitHub CLI is available."""
        try:
            subprocess.run(["gh", "--version"], capture_output=True, check=True)
            return True
        except (subprocess.CalledProcessError, FileNotFoundError):
            return False


class PluginManager:
    """Manages tldr++ plugins."""
    
    def __init__(self) -> None:
        """Initialize plugin manager."""
        self.plugins = {}
    
    def register_plugin(self, plugin: Plugin) -> None:
        """Register a plugin."""
        self.plugins[plugin.name()] = plugin
    
    def execute_plugin(self, name: str, args: List[str]) -> None:
        """Execute a plugin."""
        if name not in self.plugins:
            raise ValueError(f"Plugin '{name}' not found")
        
        self.plugins[name].execute(args)
    
    def list_plugins(self) -> List[Plugin]:
        """List all registered plugins."""
        return list(self.plugins.values())
    
    def interactive_mode(self) -> None:
        """Run the plugin in interactive mode."""
        print("tldr++ Plugin System")
        print("Type 'help' for available commands, 'exit' to quit")
        
        while True:
            try:
                command = input("tldrpp plugin> ").strip()
                
                if not command:
                    continue
                
                if command in ("exit", "quit"):
                    break
                
                if command == "help":
                    self._show_help()
                    continue
                
                if command == "list":
                    self._list_plugins()
                    continue
                
                # Parse command
                parts = command.split()
                if len(parts) < 2:
                    print("Usage: <plugin> <command> [args...]")
                    continue
                
                plugin_name = parts[0]
                args = parts[1:]
                
                self.execute_plugin(plugin_name, args)
            
            except KeyboardInterrupt:
                print("\nExiting...")
                break
            except Exception as e:
                print(f"Error: {e}")
    
    def _show_help(self) -> None:
        """Show help information."""
        print("tldr++ Plugin System")
        print()
        print("Available commands:")
        print("  help                    Show this help")
        print("  list                    List available plugins")
        print("  <plugin> <command>     Execute plugin command")
        print("  exit/quit              Exit plugin mode")
        print()
        print("Available plugins:")
        for plugin in self.plugins.values():
            print(f"  {plugin.name():<10} {plugin.description()}")
    
    def _list_plugins(self) -> None:
        """List all plugins."""
        print("Available plugins:")
        for plugin in self.plugins.values():
            print(f"  {plugin.name():<10} {plugin.description()}")