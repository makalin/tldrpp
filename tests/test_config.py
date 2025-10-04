"""Tests for configuration management."""

import tempfile
from pathlib import Path
from unittest.mock import patch

import pytest

from tldrpp.config import Config, Keymap


class TestKeymap:
    """Test Keymap class."""
    
    def test_keymap_creation(self) -> None:
        """Test keymap creation."""
        keymap = Keymap()
        assert keymap.run == "ctrl+enter"
        assert keymap.copy == "y"
        assert keymap.paste == "p"
    
    def test_keymap_custom(self) -> None:
        """Test custom keymap creation."""
        keymap = Keymap(run="enter", copy="c", paste="v")
        assert keymap.run == "enter"
        assert keymap.copy == "c"
        assert keymap.paste == "v"


class TestConfig:
    """Test Config class."""
    
    def test_config_creation(self) -> None:
        """Test config creation."""
        config = Config()
        assert config.theme == "dark"
        assert config.platforms == ["common", "linux"]
        assert config.confirm_destructive is True
        assert config.clipboard is True
        assert config.pager == "less -R"
        assert isinstance(config.keymap, Keymap)
        assert config.cache_ttl_hours == 72
        assert config.dev_mode is False
    
    def test_config_custom(self) -> None:
        """Test custom config creation."""
        keymap = Keymap(run="enter", copy="c", paste="v")
        config = Config(
            theme="light",
            platforms=["linux", "osx"],
            confirm_destructive=False,
            clipboard=False,
            pager="more",
            keymap=keymap,
            cache_ttl_hours=24,
            dev_mode=True,
        )
        
        assert config.theme == "light"
        assert config.platforms == ["linux", "osx"]
        assert config.confirm_destructive is False
        assert config.clipboard is False
        assert config.pager == "more"
        assert config.keymap.run == "enter"
        assert config.keymap.copy == "c"
        assert config.keymap.paste == "v"
        assert config.cache_ttl_hours == 24
        assert config.dev_mode is True
    
    @patch('tldrpp.config.os.path.expanduser')
    def test_get_config_file(self, mock_expanduser: Mock) -> None:
        """Test getting config file path."""
        mock_expanduser.return_value = "/home/user"
        
        config_file = Config._get_config_file()
        expected = Path("/home/user") / ".config" / "tldrpp" / "config.yml"
        assert config_file == expected
    
    @patch('tldrpp.config.os.path.expanduser')
    def test_get_default_cache_dir(self, mock_expanduser: Mock) -> None:
        """Test getting default cache directory."""
        mock_expanduser.return_value = "/home/user"
        
        cache_dir = Config._get_default_cache_dir()
        expected = str(Path("/home/user") / ".cache" / "tldrpp" / "pages")
        assert cache_dir == expected
    
    def test_save_and_load_config(self) -> None:
        """Test saving and loading config."""
        with tempfile.TemporaryDirectory() as temp_dir:
            # Create a temporary config file
            config_file = Path(temp_dir) / "config.yml"
            
            # Save config
            config = Config(theme="light", platforms=["linux"])
            with patch.object(Config, '_get_config_file', return_value=config_file):
                config.save()
            
            # Load config
            with patch.object(Config, '_get_config_file', return_value=config_file):
                loaded_config = Config.load()
            
            assert loaded_config.theme == "light"
            assert loaded_config.platforms == ["linux"]
            assert loaded_config.confirm_destructive is True  # Default value
            assert loaded_config.clipboard is True  # Default value
            assert loaded_config.pager == "less -R"  # Default value
            assert loaded_config.cache_ttl_hours == 72  # Default value
            assert loaded_config.dev_mode is False  # Default value
    
    def test_load_config_with_missing_file(self) -> None:
        """Test loading config when file doesn't exist."""
        with tempfile.TemporaryDirectory() as temp_dir:
            config_file = Path(temp_dir) / "nonexistent.yml"
            
            with patch.object(Config, '_get_config_file', return_value=config_file):
                config = Config.load()
            
            # Should return default config
            assert config.theme == "dark"
            assert config.platforms == ["common", "linux"]
            assert config.confirm_destructive is True
            assert config.clipboard is True
            assert config.pager == "less -R"
            assert config.cache_ttl_hours == 72
            assert config.dev_mode is False