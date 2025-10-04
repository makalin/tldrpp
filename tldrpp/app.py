"""Main application logic for tldr++."""

import os
import subprocess
import sys
from typing import Dict, List, Optional

from tldrpp.cache import CacheManager
from tldrpp.config import Config
from tldrpp.tui import TUIApp


class App:
    """Main application class."""
    
    def __init__(self, config: Optional[Config] = None) -> None:
        """Initialize the application."""
        self.config = config or Config.load()
        self.cache = CacheManager(self.config.cache_dir)
    
    def initialize(self) -> None:
        """Initialize tldr++ by downloading page index."""
        self.cache.initialize()
    
    def update_cache(self) -> None:
        """Update tldr pages cache."""
        self.cache.update()
    
    def run_tui(self, search_query: str = "") -> None:
        """Run the terminal user interface."""
        # Ensure cache is initialized
        if not self.cache.is_initialized():
            self.cache.initialize()
        
        app = TUIApp(self.config, self.cache)
        app.run(search_query)
    
    def render_command(self, command: str, variables: Dict[str, str]) -> str:
        """Render a command with placeholders filled."""
        page = self.cache.find_page(command)
        example = page.find_best_example(command)
        if not example:
            raise ValueError(f"No suitable example found for command: {command}")
        
        return example.render(variables)
    
    def execute_command(self, command: str, variables: Dict[str, str]) -> None:
        """Execute a command with placeholders filled."""
        page = self.cache.find_page(command)
        example = page.find_best_example(command)
        if not example:
            raise ValueError(f"No suitable example found for command: {command}")
        
        rendered = example.render(variables)
        
        # Check if command is destructive
        if self._is_destructive_command(rendered) and self.config.confirm_destructive:
            print(f"This command appears destructive: {rendered}")
            response = input("Are you sure you want to execute it? (y/N): ")
            if response.lower() not in ("y", "yes"):
                print("Command cancelled.")
                return
        
        # Execute the command
        try:
            subprocess.run(rendered, shell=True, check=True)
        except subprocess.CalledProcessError as e:
            print(f"Command failed with exit code {e.returncode}")
            sys.exit(e.returncode)
        
        # Log the execution
        self._log_execution(rendered)
    
    def submit_to_tldr(self) -> None:
        """Submit current example to tldr-pages."""
        print("Plugin system initialized. Use 'tldrpp plugin submit init' to start a submission.")
    
    def _is_destructive_command(self, command: str) -> bool:
        """Check if a command is potentially destructive."""
        destructive_verbs = [
            "rm", "rmdir", "del", "erase",
            "dd", "mkfs", "fdisk", "parted",
            "iptables", "ufw", "firewall-cmd",
            "chmod", "chown", "chattr",
            "kill", "killall", "pkill",
            "shutdown", "reboot", "halt",
            "mv", "move", "rename",
            "cp", "copy", "xcopy",
            "tar", "zip", "unzip",
            "git", "svn", "hg",
        ]
        
        command_lower = command.lower()
        for verb in destructive_verbs:
            if command_lower.startswith(f"{verb} ") or command_lower == verb:
                return True
        return False
    
    def _log_execution(self, command: str) -> None:
        """Log command execution to audit log."""
        try:
            log_dir = os.path.join(self.config.cache_dir, "..")
            os.makedirs(log_dir, exist_ok=True)
            
            log_file = os.path.join(log_dir, "exec.log")
            with open(log_file, "a") as f:
                f.write(f"{command}\n")
        except Exception:
            # Don't fail if logging fails
            pass