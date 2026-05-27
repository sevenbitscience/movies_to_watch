#!/usr/bin/env python3
import argparse
import json
import sys
from urllib.request import Request, urlopen
from urllib.error import HTTPError, URLError

# Change this if your server runs on a different port (e.g., 8080)
BASE_URL = "http://localhost:5000"

def make_request(url, method="GET", data=None):
    """Helper function to handle HTTP requests using standard library."""
    req = Request(url, method=method)
    req.add_header("Content-Type", "application/json")
    
    encoded_data = None
    if data:
        encoded_data = json.dumps(data).encode("utf-8")

    try:
        with urlopen(req, data=encoded_data) as response:
            if response.status == 204:
                return None
            return json.loads(response.read().decode("utf-8"))
    except HTTPError as e:
        print(f"[-] Error ({e.code}):", e.reason)
        try:
            # Try to print the server's custom JSON error message if it exists
            error_body = json.loads(e.read().decode("utf-8"))
            print(f"    Message: {error_body.get('message', 'No details provided.')}")
        except Exception:
            pass
        sys.exit(1)
    except URLError as e:
        print(f"[-] Could not connect to server at {BASE_URL}. Is it running?")
        print(f"    Reason: {e.reason}")
        sys.exit(1)

def list_movies(status=None):
    url = f"{BASE_URL}/movies"
    if status:
        url += f"?status={status}"
    
    movies = make_request(url)
    if not movies:
        print("[*] Your watchlist is empty!")
        return

    print(f"\n{'ID':<4} | {'Title':<30} | {'Year':<6} | {'Status':<10} | {'Genres'}")
    print("-" * 80)
    for m in movies:
        # Handle if genres come back as a list or a string representation
        genres = m.get('genre', [])
        if isinstance(genres, str):
            try:
                genres = json.loads(genres)
            except json.JSONDecodeError:
                pass
        genre_str = ", ".join(genres) if isinstance(genres, list) else str(genres)
        
        print(f"{m.get('id'):<4} | {m.get('title')[:30]:<30} | {m.get('release_year'):<6} | {m.get('status'):<10} | {genre_str}")
    print()

def search_movie(query):
    print(f"[*] Searching TMDB for '{query}'...")
    results = make_request(f"{BASE_URL}/movies/search", method="POST", data={"query": query})
    
    if not results:
        print("[-] No movies found matching that query.")
        return

    print("\nSearch Results (Use the TMDB ID to add a movie):")
    print(f"{'TMDB ID':<10} | {'Title':<35} | {'Year':<6} | {'Genres'}")
    print("-" * 80)
    for r in results:
        genres = ", ".join(r.get('genre', []))
        print(f"{r.get('tmdb_id'):<10} | {r.get('title')[:35]:<35} | {r.get('release_year'):<6} | {genres}")
    print()

def add_movie(tmdb_id, title, release_year, genre_list):
    # Expecting comma separated genres from CLI, parse them into a true list
    genres = [g.strip() for g in genre_list.split(",")] if genre_list else []
    
    payload = {
        "tmdb_id": int(tmdb_id),
        "title": title,
        "release_year": int(release_year) if release_year else None,
        "genre": genres
    }
    
    print(f"[*] Adding '{title}' to database...")
    response = make_request(f"{BASE_URL}/movies", method="POST", data=payload)
    print("[+] Movie successfully added to watchlist!")

def watch_movie(movie_id):
    print(f"[*] Marking movie ID {movie_id} as watched...")
    make_request(f"{BASE_URL}/movies/{movie_id}", method="PATCH", data={"status": "watched"})
    print(f"[+] Movie {movie_id} updated.")

def main():
    parser = argparse.ArgumentParser(description="Watchlist Sync Server CLI Client")
    subparsers = parser.add_subparsers(dest="command", help="Commands")

    # List command
    list_parser = subparsers.add_parser("list", help="List your saved movies")
    list_parser.add_argument("--status", choices=["pending", "watched"], help="Filter by status")

    # Search command
    search_parser = subparsers.add_parser("search", help="Search TMDB for a movie")
    search_parser.add_argument("query", type=str, help="Movie title to search for")

    # Add command
    add_parser = subparsers.add_parser("add", help="Add a movie manually from search results")
    add_parser.add_argument("--id", required=True, help="The TMDB ID")
    add_parser.add_argument("--title", required=True, help="Movie Title")
    add_parser.add_argument("--year", type=int, help="Release Year")
    add_parser.add_argument("--genres", help="Comma-separated genres (e.g. 'Sci-Fi,Action')")

    # Watch command
    watch_parser = subparsers.add_parser("watch", help="Mark a movie as watched by its local database ID")
    watch_parser.add_argument("id", type=int, help="Local Database Movie ID")

    args = parser.parse_args()

    if args.command == "list":
        list_movies(status=args.status)
    elif args.command == "search":
        search_movie(args.query)
    elif args.command == "add":
        add_movie(args.id, args.title, args.year, args.genres)
    elif args.command == "watch":
        watch_movie(args.id)
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
