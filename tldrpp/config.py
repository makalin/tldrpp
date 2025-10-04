"""Configuration management for tldr++."""

import os
from pathlib import Path
from typing import List

import yaml


class Config:
    """Application configuration."""
    
    def __init__(
        self,
        theme: str = "dark",
        platforms: List[str] = None,
        confirm_destructive: bool = True,
        clipboard: bool = True,
        pager: str = "less -R",
        keymap: "Keymap" = None,
        cache_ttl_hours: int = 72,
        cache_dir: str = None,
        dev_mode: bool = False,
    ) -> None:
        """Initialize configuration."""
        self.theme = theme
        self.platforms = platforms or ["common", "linux"]
        self.confirm_destructive = confirm_destructive
        self.clipboard = clipboard
        self.pager = pager
        self.keymap = keymap or Keymap()
        self.cache_ttl_hours = cache_ttl_hours
        self.cache_dir = cache_dir or self._get_default_cache_dir()
        self.dev_mode = dev_mode
    
    @classmethod
    def load(cls) -> "Config":
        """Load configuration from file or return default."""
        config_file = cls._get_config_file()
        
        if config_file.exists():
            with open(config_file) as f:
                data = yaml.safe_load(f) or {}
        else:
            data = {}
            # Create default config file
            cls._create_default_config(config_file)
        
        return cls(
            theme=data.get("theme", "dark"),
            platforms=data.get("platforms", ["common", "linux"]),
            confirm_destructive=data.get("confirm_destructive", True),
            clipboard=data.get("clipboard", True),
            pager=data.get("pager", "less -R"),
            keymap=Keymap(**data.get("keymap", {})),
            cache_ttl_hours=data.get("cache_ttl_hours", 72),
            cache_dir=data.get("cache_dir", cls._get_default_cache_dir()),
            dev_mode=data.get("dev_mode", False),
        )
    
    def save(self) -> None:
        """Save configuration to file."""
        config_file = self._get_config_file()
        config_file.parent.mkdir(parents=True, exist_ok=True)
        
        data = {
            "theme": self.theme,
            "platforms": self.platforms,
            "confirm_destructive": self.confirm_destructive,
            "clipboard": self.clipboard,
            "pager": self.pager,
            "keymap": {
                "run": self.keymap.run,
                "copy": self.keymap.copy,
                "paste": self.keymap.paste,
            },
            "cache_ttl_hours": self.cache_ttl_hours,
            "cache_dir": self.cache_dir,
            "dev_mode": self.dev_mode,
        }
        
        with open(config_file, "w") as f:
            yaml.dump(data, f, default_flow_style=False)
    
    @staticmethod
    def _get_config_file() -> Path:
        """Get the configuration file path."""
        if home_dir := os.path.expanduser("~"):
            return Path(home_dir) / ".config" / "tldrpp" / "config.yml"
        return Path(".config") / "tldrpp" / "config.yml"
    
    @staticmethod
    def _get_default_cache_dir() -> str:
        """Get the default cache directory."""
        if home_dir := os.path.expanduser("~"):
            return str(Path(home_dir) / ".cache" / "tldrpp" / "pages")
        return str(Path(".cache") / "tldrpp" / "pages")
    
    @staticmethod
    def _create_default_config(config_file: Path) -> None:
        """Create a default configuration file."""
        config_file.parent.mkdir(parents=True, exist_ok=True)
        
        default_config = Config()
        default_config.save()


class Keymap:
    """Keyboard shortcuts configuration."""
    
    def __init__(
        self,
        run: str = "ctrl+enter",
        copy: str = "y",
        paste: str = "p",
    ) -> None:
        """Initialize keymap."""
        self.run = run
        self.copy = copy
        self.paste = paste