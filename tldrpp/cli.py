"""Command line interface for tldr++."""

import sys
from typing import Optional

import click

from tldrpp.app import App
from tldrpp.config import Config


@click.group()
@click.version_option(version="0.1.0")
@click.option(
    "--platform", "-p", 
    help="Platform filter (common, linux, osx, sunos, windows, android)"
)
@click.option(
    "--theme", "-t", 
    default="dark", 
    help="Theme (light, dark, solarized)"
)
@click.option(
    "--dev", "-d", 
    is_flag=True, 
    help="Development mode"
)
@click.pass_context
def cli(ctx: click.Context, platform: Optional[str], theme: str, dev: bool) -> None:
    """tldr++ - Interactive cheat-sheets with fuzzy search and inline editing."""
    ctx.ensure_object(dict)
    ctx.obj["platform"] = platform
    ctx.obj["theme"] = theme
    ctx.obj["dev"] = dev


@cli.command()
def init() -> None:
    """Initialize tldr++ by downloading page index."""
    try:
        app = App()
        app.initialize()
        click.echo("tldr++ initialized successfully!")
    except Exception as e:
        click.echo(f"Error initializing tldr++: {e}", err=True)
        sys.exit(1)


@cli.command()
def update() -> None:
    """Update tldr pages cache."""
    try:
        app = App()
        app.update_cache()
        click.echo("Cache updated successfully!")
    except Exception as e:
        click.echo(f"Error updating cache: {e}", err=True)
        sys.exit(1)


@cli.command()
@click.argument("command")
@click.option(
    "--vars", 
    help="Variables to substitute in placeholders (format: key=value,key2=value2)"
)
def render(command: str, vars: Optional[str]) -> None:
    """Render command with placeholders filled."""
    try:
        app = App()
        variables = {}
        if vars:
            for pair in vars.split(","):
                if "=" in pair:
                    key, value = pair.split("=", 1)
                    variables[key.strip()] = value.strip()
        
        result = app.render_command(command, variables)
        click.echo(result)
    except Exception as e:
        click.echo(f"Error rendering command: {e}", err=True)
        sys.exit(1)


@cli.command()
@click.argument("command")
@click.option(
    "--vars", 
    help="Variables to substitute in placeholders (format: key=value,key2=value2)"
)
def exec(command: str, vars: Optional[str]) -> None:
    """Execute command with placeholders filled."""
    try:
        app = App()
        variables = {}
        if vars:
            for pair in vars.split(","):
                if "=" in pair:
                    key, value = pair.split("=", 1)
                    variables[key.strip()] = value.strip()
        
        app.execute_command(command, variables)
    except Exception as e:
        click.echo(f"Error executing command: {e}", err=True)
        sys.exit(1)


@cli.group()
def plugin() -> None:
    """Plugin commands."""
    pass


@plugin.command()
def submit() -> None:
    """Submit current example to tldr-pages."""
    try:
        app = App()
        app.submit_to_tldr()
    except Exception as e:
        click.echo(f"Error submitting to tldr: {e}", err=True)
        sys.exit(1)


@cli.command()
@click.argument("search_query", required=False)
@click.pass_context
def run(ctx: click.Context, search_query: Optional[str]) -> None:
    """Run the tldr++ TUI (default command)."""
    try:
        config = Config.load()
        
        # Override config with command line flags
        if ctx.obj["platform"]:
            config.platforms = [ctx.obj["platform"]]
        if ctx.obj["theme"]:
            config.theme = ctx.obj["theme"]
        if ctx.obj["dev"]:
            config.dev_mode = True
        
        app = App(config)
        app.run_tui(search_query or "")
    except Exception as e:
        click.echo(f"Error running tldr++: {e}", err=True)
        sys.exit(1)


def main() -> None:
    """Main entry point."""
    cli()