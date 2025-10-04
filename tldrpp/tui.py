"""Terminal user interface for tldr++."""

from typing import List, Optional

from textual.app import App as TextualApp, ComposeResult
from textual.containers import Container, Horizontal, Vertical
from textual.reactive import reactive
from textual.widgets import (
    Button,
    DataTable,
    Footer,
    Header,
    Input,
    Label,
    ListView,
    ListItem,
    Static,
    TabbedContent,
    Tab,
    TextArea,
)

from tldrpp.cache import CacheManager, Page
from tldrpp.config import Config


class SearchWidget(Static):
    """Search input widget."""
    
    def __init__(self, placeholder: str = "Search commands...") -> None:
        """Initialize search widget."""
        super().__init__()
        self.placeholder = placeholder
    
    def compose(self) -> ComposeResult:
        """Compose the widget."""
        yield Input(placeholder=self.placeholder, id="search_input")


class PagesWidget(Static):
    """Pages list widget."""
    
    def __init__(self, pages: List[Page]) -> None:
        """Initialize pages widget."""
        super().__init__()
        self.pages = pages
    
    def compose(self) -> ComposeResult:
        """Compose the widget."""
        with ListView(id="pages_list"):
            for page in self.pages:
                yield ListItem(
                    Label(f"{page.name} - {page.description} ({page.platform})"),
                    id=f"page_{page.name}"
                )


class ExamplesWidget(Static):
    """Examples widget."""
    
    def __init__(self, page: Page) -> None:
        """Initialize examples widget."""
        super().__init__()
        self.page = page
    
    def compose(self) -> ComposeResult:
        """Compose the widget."""
        with Vertical():
            yield Label(f"Examples for {self.page.name}", id="examples_title")
            
            with TabbedContent():
                for i, example in enumerate(self.page.examples):
                    with Tab(f"Example {i+1}", id=f"example_{i}"):
                        yield Label(example.description, id=f"desc_{i}")
                        yield TextArea(
                            example.command,
                            read_only=True,
                            id=f"command_{i}"
                        )


class TUIApp(TextualApp):
    """Main TUI application."""
    
    CSS = """
    Screen {
        layout: vertical;
    }
    
    Header {
        dock: top;
    }
    
    Footer {
        dock: bottom;
    }
    
    #main_container {
        layout: horizontal;
        height: 1fr;
    }
    
    #left_panel {
        width: 30%;
        border: solid $primary;
        margin: 1;
    }
    
    #right_panel {
        width: 70%;
        border: solid $primary;
        margin: 1;
    }
    
    #search_input {
        margin: 1;
    }
    
    #pages_list {
        height: 1fr;
    }
    
    #examples_title {
        text-style: bold;
        margin: 1;
    }
    
    TabbedContent {
        height: 1fr;
    }
    
    Tab {
        padding: 1;
    }
    
    TextArea {
        height: 1fr;
        margin: 1;
    }
    """
    
    def __init__(self, config: Config, cache: CacheManager) -> None:
        """Initialize TUI app."""
        super().__init__()
        self.config = config
        self.cache = cache
        self.pages: List[Page] = []
        self.selected_page: Optional[Page] = None
    
    def compose(self) -> ComposeResult:
        """Compose the app."""
        yield Header()
        
        with Container(id="main_container"):
            with Vertical(id="left_panel"):
                yield SearchWidget()
                yield PagesWidget(self.pages)
            
            with Vertical(id="right_panel"):
                yield Static("Select a page to view examples", id="examples_placeholder")
        
        yield Footer()
    
    def on_mount(self) -> None:
        """Handle app mount."""
        self.title = "tldr++ - Interactive Cheat-Sheets"
        self.sub_title = "Fuzzy search and inline editing for tldr pages"
        
        # Load initial pages
        self.load_pages("")
    
    def on_input_changed(self, event: Input.Changed) -> None:
        """Handle input changes."""
        if event.input.id == "search_input":
            self.load_pages(event.value)
    
    def on_list_view_selected(self, event: ListView.Selected) -> None:
        """Handle list selection."""
        if event.list_view.id == "pages_list":
            selected_item = event.item
            if selected_item:
                page_name = selected_item.id.replace("page_", "")
                self.selected_page = next(
                    (p for p in self.pages if p.name == page_name), None
                )
                if self.selected_page:
                    self.show_examples()
    
    def load_pages(self, query: str) -> None:
        """Load pages based on search query."""
        try:
            self.pages = self.cache.search_pages(query, self.config.platforms)
            self.refresh_pages_list()
        except Exception as e:
            self.notify(f"Error loading pages: {e}", severity="error")
    
    def refresh_pages_list(self) -> None:
        """Refresh the pages list widget."""
        pages_list = self.query_one("#pages_list", ListView)
        pages_list.clear()
        
        for page in self.pages:
            item = ListItem(
                Label(f"{page.name} - {page.description} ({page.platform})"),
                id=f"page_{page.name}"
            )
            pages_list.append(item)
    
    def show_examples(self) -> None:
        """Show examples for the selected page."""
        if not self.selected_page:
            return
        
        # Replace the right panel content
        right_panel = self.query_one("#right_panel", Vertical)
        right_panel.remove_children()
        
        examples_widget = ExamplesWidget(self.selected_page)
        right_panel.mount(examples_widget)
    
    def run(self, search_query: str = "") -> None:
        """Run the TUI."""
        self.run()