"""Cache management for tldr pages."""

import json
import os
import re
from pathlib import Path
from typing import Dict, List, Optional

import requests


class IndexEntry:
    """Entry in the tldr pages index."""
    
    def __init__(self, name: str, description: str, platform: str) -> None:
        """Initialize index entry."""
        self.name = name
        self.description = description
        self.platform = platform


class Placeholder:
    """Placeholder in a command."""
    
    def __init__(
        self,
        name: str,
        type_: str = "text",
        description: str = "",
        default: str = "",
    ) -> None:
        """Initialize placeholder."""
        self.name = name
        self.type = type_
        self.description = description
        self.default = default


class Example:
    """Command example."""
    
    def __init__(
        self,
        description: str,
        command: str,
        placeholders: List[Placeholder] = None,
    ) -> None:
        """Initialize example."""
        self.description = description
        self.command = command
        self.placeholders = placeholders or []
    
    def render(self, variables: Dict[str, str]) -> str:
        """Render command with placeholders filled."""
        command = self.command
        
        for placeholder in self.placeholders:
            value = variables.get(placeholder.name, placeholder.default)
            if not value:
                value = placeholder.name  # Use placeholder name as fallback
            
            pattern = re.compile(rf"\{\{{{re.escape(placeholder.name)}}}\}}")
            command = pattern.sub(value, command)
        
        return command


class Page:
    """tldr page."""
    
    def __init__(
        self,
        name: str,
        description: str,
        platform: str,
        examples: List[Example] = None,
        raw_content: str = "",
    ) -> None:
        """Initialize page."""
        self.name = name
        self.description = description
        self.platform = platform
        self.examples = examples or []
        self.raw_content = raw_content
    
    def find_best_example(self, query: str) -> Optional[Example]:
        """Find the best matching example for a command."""
        if not self.examples:
            return None
        
        query_lower = query.lower()
        
        # Look for exact match in description
        for example in self.examples:
            if query_lower in example.description.lower():
                return example
        
        # Look for partial match
        for example in self.examples:
            if query_lower in example.description.lower():
                return example
        
        # Return first example as fallback
        return self.examples[0]


class CacheManager:
    """Manages tldr pages caching."""
    
    def __init__(self, cache_dir: str) -> None:
        """Initialize cache manager."""
        self.cache_dir = Path(cache_dir)
        self.session = requests.Session()
        self.session.timeout = 30
    
    def initialize(self) -> None:
        """Initialize cache by downloading pages."""
        if self.is_initialized():
            return
        
        # Download pages index
        index = self._download_index()
        
        # Download all pages
        self._download_pages(index)
        
        # Save index
        self._save_index(index)
    
    def update(self) -> None:
        """Update cache."""
        self.initialize()
    
    def is_initialized(self) -> bool:
        """Check if cache is initialized."""
        index_file = self.cache_dir / "index.json"
        return index_file.exists()
    
    def find_page(self, command: str) -> Page:
        """Find a page by command name."""
        index = self._load_index()
        
        # Search for exact match first
        for entry in index:
            if entry.name == command:
                return self._load_page(entry)
        
        # Search for partial matches
        matches = []
        for entry in index:
            if command.lower() in entry.name.lower():
                matches.append(entry)
        
        if not matches:
            raise ValueError(f"Command not found: {command}")
        
        # Sort by relevance (exact prefix matches first)
        matches.sort(key=lambda x: (
            not x.name.lower().startswith(command.lower()),
            x.name.lower()
        ))
        
        return self._load_page(matches[0])
    
    def search_pages(self, query: str, platforms: List[str]) -> List[Page]:
        """Search for pages matching a query."""
        index = self._load_index()
        results = []
        query_lower = query.lower()
        
        for entry in index:
            # Filter by platform if specified
            if platforms and entry.platform not in platforms:
                continue
            
            # Check if query matches
            if (query_lower in entry.name.lower() or 
                query_lower in entry.description.lower()):
                try:
                    page = self._load_page(entry)
                    results.append(page)
                except Exception:
                    # Skip pages that can't be loaded
                    continue
        
        # Sort by relevance
        results.sort(key=lambda x: self._calculate_relevance_score(x, query))
        return results
    
    def _download_index(self) -> List[IndexEntry]:
        """Download the pages index from tldr-pages."""
        url = "https://raw.githubusercontent.com/tldr-pages/tldr/main/pages.json"
        response = self.session.get(url)
        response.raise_for_status()
        
        data = response.json()
        return [
            IndexEntry(
                name=item["name"],
                description=item["description"],
                platform=item["platform"]
            )
            for item in data
        ]
    
    def _download_pages(self, index: List[IndexEntry]) -> None:
        """Download all pages."""
        for entry in index:
            try:
                self._download_page(entry)
            except Exception as e:
                print(f"Warning: failed to download page {entry.name}: {e}")
    
    def _download_page(self, entry: IndexEntry) -> None:
        """Download a single page."""
        url = (f"https://raw.githubusercontent.com/tldr-pages/tldr/main/pages/"
               f"{entry.platform}/{entry.name}.md")
        
        response = self.session.get(url)
        response.raise_for_status()
        
        # Create platform directory
        platform_dir = self.cache_dir / entry.platform
        platform_dir.mkdir(parents=True, exist_ok=True)
        
        # Save page
        page_file = platform_dir / f"{entry.name}.md"
        with open(page_file, "w", encoding="utf-8") as f:
            f.write(response.text)
    
    def _save_index(self, index: List[IndexEntry]) -> None:
        """Save the index to disk."""
        self.cache_dir.mkdir(parents=True, exist_ok=True)
        index_file = self.cache_dir / "index.json"
        
        data = [
            {
                "name": entry.name,
                "description": entry.description,
                "platform": entry.platform
            }
            for entry in index
        ]
        
        with open(index_file, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2)
    
    def _load_index(self) -> List[IndexEntry]:
        """Load the index from disk."""
        index_file = self.cache_dir / "index.json"
        
        with open(index_file, encoding="utf-8") as f:
            data = json.load(f)
        
        return [
            IndexEntry(
                name=item["name"],
                description=item["description"],
                platform=item["platform"]
            )
            for item in data
        ]
    
    def _load_page(self, entry: IndexEntry) -> Page:
        """Load a page from disk."""
        page_file = self.cache_dir / entry.platform / f"{entry.name}.md"
        
        with open(page_file, encoding="utf-8") as f:
            content = f.read()
        
        return self._parse_page(content, entry)
    
    def _parse_page(self, content: str, entry: IndexEntry) -> Page:
        """Parse a tldr page from markdown content."""
        page = Page(
            name=entry.name,
            description=entry.description,
            platform=entry.platform,
            raw_content=content
        )
        
        lines = content.split("\n")
        current_example = None
        in_example = False
        
        for line in lines:
            line = line.strip()
            
            if line.startswith("# "):
                # Skip title
                continue
            elif line.startswith("> "):
                # Description
                page.description = line[2:]
            elif line.startswith("- "):
                # Start new example
                if current_example:
                    page.examples.append(current_example)
                
                current_example = Example(description=line[2:])
                in_example = True
            elif line.startswith("`") and line.endswith("`") and in_example:
                # Command
                command = line[1:-1]
                current_example.command = command
                current_example.placeholders = self._extract_placeholders(command)
            elif not line:
                # Empty line ends example
                in_example = False
        
        # Add last example
        if current_example:
            page.examples.append(current_example)
        
        return page
    
    def _extract_placeholders(self, command: str) -> List[Placeholder]:
        """Extract placeholders from a command string."""
        placeholders = []
        
        # Regex to find {{placeholder}} patterns
        pattern = re.compile(r"\{\{([^}]+)\}\}")
        matches = pattern.findall(command)
        
        seen = set()
        for match in matches:
            if match not in seen:
                seen.add(match)
                placeholder = Placeholder(
                    name=match,
                    type=self._infer_placeholder_type(match)
                )
                placeholders.append(placeholder)
        
        return placeholders
    
    def _infer_placeholder_type(self, name: str) -> str:
        """Infer the type of a placeholder based on its name."""
        name_lower = name.lower()
        
        if "file" in name_lower or "path" in name_lower:
            return "file"
        elif "dir" in name_lower or "directory" in name_lower:
            return "directory"
        elif "port" in name_lower:
            return "port"
        elif "num" in name_lower or "number" in name_lower or "count" in name_lower:
            return "number"
        elif "url" in name_lower or "link" in name_lower:
            return "url"
        elif "ip" in name_lower or "address" in name_lower:
            return "ip"
        elif "user" in name_lower or "username" in name_lower:
            return "username"
        elif "pass" in name_lower or "password" in name_lower:
            return "password"
        elif "email" in name_lower:
            return "email"
        else:
            return "text"
    
    def _calculate_relevance_score(self, page: Page, query: str) -> int:
        """Calculate relevance score for search results."""
        score = 0
        query_lower = query.lower()
        name_lower = page.name.lower()
        description_lower = page.description.lower()
        
        # Exact name match gets highest score
        if name_lower == query_lower:
            score += 100
        elif name_lower.startswith(query_lower):
            score += 50
        elif query_lower in name_lower:
            score += 25
        
        # Description match gets lower score
        if query_lower in description_lower:
            score += 10
        
        # Example matches get medium score
        for example in page.examples:
            if query_lower in example.description.lower():
                score += 15
        
        return score