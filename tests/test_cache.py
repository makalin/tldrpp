"""Tests for cache functionality."""

import tempfile
from pathlib import Path
from unittest.mock import Mock, patch

import pytest

from tldrpp.cache import CacheManager, Example, Page, Placeholder


class TestPlaceholder:
    """Test Placeholder class."""
    
    def test_placeholder_creation(self) -> None:
        """Test placeholder creation."""
        placeholder = Placeholder("file", "file", "Input file", "test.txt")
        assert placeholder.name == "file"
        assert placeholder.type == "file"
        assert placeholder.description == "Input file"
        assert placeholder.default == "test.txt"


class TestExample:
    """Test Example class."""
    
    def test_example_creation(self) -> None:
        """Test example creation."""
        example = Example("Extract archive", "tar -xf {{file}}")
        assert example.description == "Extract archive"
        assert example.command == "tar -xf {{file}}"
        assert len(example.placeholders) == 1
        assert example.placeholders[0].name == "file"
    
    def test_example_render(self) -> None:
        """Test example rendering with variables."""
        example = Example("Extract archive", "tar -xf {{file}}")
        variables = {"file": "archive.tar.gz"}
        result = example.render(variables)
        assert result == "tar -xf archive.tar.gz"
    
    def test_example_render_with_default(self) -> None:
        """Test example rendering with default values."""
        placeholder = Placeholder("file", "file", "Input file", "default.tar.gz")
        example = Example("Extract archive", "tar -xf {{file}}", [placeholder])
        result = example.render({})
        assert result == "tar -xf default.tar.gz"


class TestPage:
    """Test Page class."""
    
    def test_page_creation(self) -> None:
        """Test page creation."""
        examples = [
            Example("Extract archive", "tar -xf {{file}}"),
            Example("List contents", "tar -tf {{file}}"),
        ]
        page = Page("tar", "Archive utility", "linux", examples)
        assert page.name == "tar"
        assert page.description == "Archive utility"
        assert page.platform == "linux"
        assert len(page.examples) == 2
    
    def test_find_best_example(self) -> None:
        """Test finding best example."""
        examples = [
            Example("Extract archive", "tar -xf {{file}}"),
            Example("List contents", "tar -tf {{file}}"),
        ]
        page = Page("tar", "Archive utility", "linux", examples)
        
        # Should return first example for generic query
        best = page.find_best_example("tar")
        assert best is not None
        assert best.description == "Extract archive"
        
        # Should return None for empty examples
        empty_page = Page("empty", "Empty page", "linux", [])
        assert empty_page.find_best_example("query") is None


class TestCacheManager:
    """Test CacheManager class."""
    
    def test_cache_manager_creation(self) -> None:
        """Test cache manager creation."""
        with tempfile.TemporaryDirectory() as temp_dir:
            cache = CacheManager(temp_dir)
            assert cache.cache_dir == Path(temp_dir)
    
    def test_is_initialized_false(self) -> None:
        """Test is_initialized returns False for empty cache."""
        with tempfile.TemporaryDirectory() as temp_dir:
            cache = CacheManager(temp_dir)
            assert not cache.is_initialized()
    
    def test_is_initialized_true(self) -> None:
        """Test is_initialized returns True when index exists."""
        with tempfile.TemporaryDirectory() as temp_dir:
            cache = CacheManager(temp_dir)
            cache.cache_dir.mkdir(parents=True, exist_ok=True)
            (cache.cache_dir / "index.json").touch()
            assert cache.is_initialized()
    
    def test_extract_placeholders(self) -> None:
        """Test placeholder extraction."""
        with tempfile.TemporaryDirectory() as temp_dir:
            cache = CacheManager(temp_dir)
            
            # Test single placeholder
            placeholders = cache._extract_placeholders("tar -xf {{file}}")
            assert len(placeholders) == 1
            assert placeholders[0].name == "file"
            
            # Test multiple placeholders
            placeholders = cache._extract_placeholders("cp {{src}} {{dest}}")
            assert len(placeholders) == 2
            assert placeholders[0].name == "src"
            assert placeholders[1].name == "dest"
            
            # Test no placeholders
            placeholders = cache._extract_placeholders("ls -la")
            assert len(placeholders) == 0
    
    def test_infer_placeholder_type(self) -> None:
        """Test placeholder type inference."""
        with tempfile.TemporaryDirectory() as temp_dir:
            cache = CacheManager(temp_dir)
            
            assert cache._infer_placeholder_type("file") == "file"
            assert cache._infer_placeholder_type("directory") == "directory"
            assert cache._infer_placeholder_type("port") == "port"
            assert cache._infer_placeholder_type("number") == "number"
            assert cache._infer_placeholder_type("url") == "url"
            assert cache._infer_placeholder_type("ip") == "ip"
            assert cache._infer_placeholder_type("username") == "username"
            assert cache._infer_placeholder_type("password") == "password"
            assert cache._infer_placeholder_type("email") == "email"
            assert cache._infer_placeholder_type("unknown") == "text"
    
    def test_parse_page(self) -> None:
        """Test page parsing."""
        with tempfile.TemporaryDirectory() as temp_dir:
            cache = CacheManager(temp_dir)
            
            content = """# tar

> Archive utility.

- Extract archive:
  `tar -xf {{file}}`

- List contents:
  `tar -tf {{file}}`
"""
            
            from tldrpp.cache import IndexEntry
            entry = IndexEntry("tar", "Archive utility", "linux")
            page = cache._parse_page(content, entry)
            
            assert page.name == "tar"
            assert page.description == "Archive utility"
            assert page.platform == "linux"
            assert len(page.examples) == 2
            assert page.examples[0].description == "Extract archive"
            assert page.examples[0].command == "tar -xf {{file}}"
            assert page.examples[1].description == "List contents"
            assert page.examples[1].command == "tar -tf {{file}}"
    
    @patch('tldrpp.cache.requests.Session.get')
    def test_download_index(self, mock_get: Mock) -> None:
        """Test downloading index."""
        with tempfile.TemporaryDirectory() as temp_dir:
            cache = CacheManager(temp_dir)
            
            # Mock response
            mock_response = Mock()
            mock_response.json.return_value = [
                {"name": "tar", "description": "Archive utility", "platform": "linux"},
                {"name": "ls", "description": "List files", "platform": "common"},
            ]
            mock_response.raise_for_status.return_value = None
            mock_get.return_value = mock_response
            
            index = cache._download_index()
            
            assert len(index) == 2
            assert index[0].name == "tar"
            assert index[0].description == "Archive utility"
            assert index[0].platform == "linux"
            assert index[1].name == "ls"
            assert index[1].description == "List files"
            assert index[1].platform == "common"
    
    def test_save_and_load_index(self) -> None:
        """Test saving and loading index."""
        with tempfile.TemporaryDirectory() as temp_dir:
            cache = CacheManager(temp_dir)
            
            from tldrpp.cache import IndexEntry
            index = [
                IndexEntry("tar", "Archive utility", "linux"),
                IndexEntry("ls", "List files", "common"),
            ]
            
            cache._save_index(index)
            loaded_index = cache._load_index()
            
            assert len(loaded_index) == 2
            assert loaded_index[0].name == "tar"
            assert loaded_index[0].description == "Archive utility"
            assert loaded_index[0].platform == "linux"
            assert loaded_index[1].name == "ls"
            assert loaded_index[1].description == "List files"
            assert loaded_index[1].platform == "common"