"""Tests for application logic."""

import tempfile
from pathlib import Path
from unittest.mock import Mock, patch

import pytest

from tldrpp.app import App
from tldrpp.cache import CacheManager, Example, Page
from tldrpp.config import Config


class TestApp:
    """Test App class."""
    
    def test_app_creation(self) -> None:
        """Test app creation."""
        with tempfile.TemporaryDirectory() as temp_dir:
            config = Config(cache_dir=temp_dir)
            app = App(config)
            assert app.config == config
            assert isinstance(app.cache, CacheManager)
    
    def test_app_creation_with_default_config(self) -> None:
        """Test app creation with default config."""
        app = App()
        assert isinstance(app.config, Config)
        assert isinstance(app.cache, CacheManager)
    
    @patch('tldrpp.app.CacheManager.initialize')
    def test_initialize(self, mock_initialize: Mock) -> None:
        """Test app initialization."""
        app = App()
        app.initialize()
        mock_initialize.assert_called_once()
    
    @patch('tldrpp.app.CacheManager.update')
    def test_update_cache(self, mock_update: Mock) -> None:
        """Test cache update."""
        app = App()
        app.update_cache()
        mock_update.assert_called_once()
    
    @patch('tldrpp.app.CacheManager.is_initialized')
    @patch('tldrpp.app.CacheManager.initialize')
    @patch('tldrpp.app.TUIApp')
    def test_run_tui(self, mock_tui_app: Mock, mock_initialize: Mock, mock_is_initialized: Mock) -> None:
        """Test running TUI."""
        mock_is_initialized.return_value = False
        mock_tui_instance = Mock()
        mock_tui_app.return_value = mock_tui_instance
        
        app = App()
        app.run_tui("test query")
        
        mock_initialize.assert_called_once()
        mock_tui_app.assert_called_once_with(app.config, app.cache)
        mock_tui_instance.run.assert_called_once_with("test query")
    
    @patch('tldrpp.app.CacheManager.is_initialized')
    @patch('tldrpp.app.CacheManager.find_page')
    def test_render_command(self, mock_find_page: Mock, mock_is_initialized: Mock) -> None:
        """Test rendering command."""
        mock_is_initialized.return_value = True
        
        # Create mock page and example
        example = Example("Extract archive", "tar -xf {{file}}")
        page = Page("tar", "Archive utility", "linux", [example])
        mock_find_page.return_value = page
        
        app = App()
        result = app.render_command("tar", {"file": "archive.tar.gz"})
        
        assert result == "tar -xf archive.tar.gz"
        mock_find_page.assert_called_once_with("tar")
    
    @patch('tldrpp.app.CacheManager.is_initialized')
    @patch('tldrpp.app.CacheManager.find_page')
    def test_render_command_no_example(self, mock_find_page: Mock, mock_is_initialized: Mock) -> None:
        """Test rendering command with no suitable example."""
        mock_is_initialized.return_value = True
        
        # Create mock page with no examples
        page = Page("tar", "Archive utility", "linux", [])
        mock_find_page.return_value = page
        
        app = App()
        
        with pytest.raises(ValueError, match="No suitable example found"):
            app.render_command("tar", {})
    
    @patch('tldrpp.app.CacheManager.is_initialized')
    @patch('tldrpp.app.CacheManager.find_page')
    @patch('tldrpp.app.subprocess.run')
    def test_execute_command(self, mock_run: Mock, mock_find_page: Mock, mock_is_initialized: Mock) -> None:
        """Test executing command."""
        mock_is_initialized.return_value = True
        
        # Create mock page and example
        example = Example("Extract archive", "tar -xf {{file}}")
        page = Page("tar", "Archive utility", "linux", [example])
        mock_find_page.return_value = page
        
        app = App()
        app.execute_command("tar", {"file": "archive.tar.gz"})
        
        mock_run.assert_called_once_with("tar -xf archive.tar.gz", shell=True, check=True)
    
    @patch('tldrpp.app.CacheManager.is_initialized')
    @patch('tldrpp.app.CacheManager.find_page')
    @patch('tldrpp.app.subprocess.run')
    def test_execute_command_destructive(self, mock_run: Mock, mock_find_page: Mock, mock_is_initialized: Mock) -> None:
        """Test executing destructive command."""
        mock_is_initialized.return_value = True
        
        # Create mock page and example with destructive command
        example = Example("Remove file", "rm {{file}}")
        page = Page("rm", "Remove files", "linux", [example])
        mock_find_page.return_value = page
        
        app = App()
        
        with patch('builtins.input', return_value='n'):
            app.execute_command("rm", {"file": "test.txt"})
        
        # Should not execute the command
        mock_run.assert_not_called()
    
    def test_is_destructive_command(self) -> None:
        """Test destructive command detection."""
        app = App()
        
        # Test destructive commands
        assert app._is_destructive_command("rm file.txt")
        assert app._is_destructive_command("dd if=/dev/zero of=file")
        assert app._is_destructive_command("chmod 777 file")
        assert app._is_destructive_command("kill 1234")
        assert app._is_destructive_command("shutdown now")
        
        # Test non-destructive commands
        assert not app._is_destructive_command("ls -la")
        assert not app._is_destructive_command("cat file.txt")
        assert not app._is_destructive_command("echo hello")
        assert not app._is_destructive_command("pwd")
    
    def test_submit_to_tldr(self) -> None:
        """Test submitting to tldr."""
        app = App()
        
        # Should not raise an exception
        app.submit_to_tldr()
    
    @patch('tldrpp.app.os.makedirs')
    @patch('builtins.open', create=True)
    def test_log_execution(self, mock_open: Mock, mock_makedirs: Mock) -> None:
        """Test logging execution."""
        app = App()
        
        mock_file = Mock()
        mock_open.return_value.__enter__.return_value = mock_file
        
        app._log_execution("test command")
        
        mock_makedirs.assert_called_once()
        mock_file.write.assert_called_once_with("test command\n")