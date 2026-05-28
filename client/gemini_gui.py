#!/usr/bin/env python3
import json
import tkinter as tk
from tkinter import ttk, messagebox
from urllib.request import Request, urlopen
from urllib.error import HTTPError, URLError

# Configuration
BASE_URL = "http://localhost:5000"

class WatchlistApp(tk.Tk):
    def __init__(self):
        super().__init__()
        self.title("Watchlist Sync Center")
        self.geometry("900x550")
        self.minsize(800, 450)
        
        # Configure styles
        self.style = ttk.Style()
        self.style.theme_use("clam")
        
        # Main layout splitting Watchlist (Left) and TMDB Search (Right)
        main_paned = ttk.Panedwindow(self, orient=tk.HORIZONTAL)
        main_paned.pack(fill=tk.BOTH, expand=True, padx=10, pady=10)
        
        # --- LEFT FRAME: CURRENT WATCHLIST ---
        left_frame = ttk.LabelFrame(main_paned, text=" My Watchlist ", padding=10)
        main_paned.add(left_frame, weight=3)
        
        # Watchlist Filter Buttons
        filter_frame = ttk.Frame(left_frame)
        filter_frame.pack(fill=tk.X, pady=(0, 5))
        ttk.Button(filter_frame, text="Refresh All", command=lambda: self.load_watchlist()).pack(side=tk.LEFT, padx=2)
        ttk.Button(filter_frame, text="Pending Only", command=lambda: self.load_watchlist("pending")).pack(side=tk.LEFT, padx=2)
        ttk.Button(filter_frame, text="Watched Only", command=lambda: self.load_watchlist("watched")).pack(side=tk.LEFT, padx=2)
        
        # Watchlist Treeview Table
        self.watchlist_tree = ttk.Treeview(left_frame, columns=("id", "title", "year", "genres", "status"), show="headings")
        self.watchlist_tree.heading("id", text="ID")
        self.watchlist_tree.heading("title", text="Title")
        self.watchlist_tree.heading("year", text="Year")
        self.watchlist_tree.heading("genres", text="Genres")
        self.watchlist_tree.heading("status", text="Status")
        
        self.watchlist_tree.column("id", width=40, anchor=tk.CENTER)
        self.watchlist_tree.column("title", width=180, anchor=tk.W)
        self.watchlist_tree.column("year", width=60, anchor=tk.CENTER)
        self.watchlist_tree.column("genres", width=150, anchor=tk.W)
        self.watchlist_tree.column("status", width=80, anchor=tk.CENTER)
        self.watchlist_tree.pack(fill=tk.BOTH, expand=True)
        
        # Watchlist Action Button
        ttk.Button(left_frame, text="Toggle Watched / Pending Status", command=self.toggle_movie_status).pack(fill=tk.X, pady=(5, 0))
        
        # --- RIGHT FRAME: SEARCH & ADD ---
        right_frame = ttk.LabelFrame(main_paned, text=" Find & Add Movies (TMDB) ", padding=10)
        main_paned.add(right_frame, weight=2)
        
        # Search Inputs
        search_bar_frame = ttk.Frame(right_frame)
        search_bar_frame.pack(fill=tk.X, pady=(0, 5))
        self.search_entry = ttk.Entry(search_bar_frame)
        self.search_entry.pack(side=tk.LEFT, fill=tk.X, expand=True, padx=(0, 5))
        self.search_entry.bind("<Return>", lambda e: self.search_tmdb())
        
        ttk.Button(search_bar_frame, text="Search", command=self.search_tmdb).pack(side=tk.RIGHT)
        
        # Search Results Treeview Table
        self.search_tree = ttk.Treeview(right_frame, columns=("tmdb_id", "title", "year", "genres"), show="headings")
        self.search_tree.heading("tmdb_id", text="TMDB ID")
        self.search_tree.heading("title", text="Title")
        self.search_tree.heading("year", text="Year")
        self.search_tree.heading("genres", text="Genres")
        
        self.search_tree.column("tmdb_id", width=70, anchor=tk.CENTER)
        self.search_tree.column("title", width=140, anchor=tk.W)
        self.search_tree.column("year", width=50, anchor=tk.CENTER)
        self.search_tree.column("genres", width=120, anchor=tk.W)
        self.search_tree.pack(fill=tk.BOTH, expand=True)
        
        # Add Selected Movie Button
        ttk.Button(right_frame, text="Add Selected to Watchlist", command=self.add_selected_movie).pack(fill=tk.X, pady=(5, 0))
        
        # Initial load of tracking database
        self.load_watchlist()

    def make_network_request(self, endpoint, method="GET", data=None):
        """Standard HTTP execution matching specified empty response guidelines."""
        url = f"{BASE_URL}{endpoint}"
        req = Request(url, method=method)
        req.add_header("Content-Type", "application/json")
        
        encoded_payload = json.dumps(data).encode("utf-8") if data else None
        
        try:
            with urlopen(req, data=encoded_payload) as response:
                raw_bytes = response.read()
                if not raw_bytes or raw_bytes.strip() == b"":
                    return None
                return json.loads(raw_bytes.decode("utf-8"))
        except HTTPError as e:
            msg = e.reason
            try:
                error_body = json.loads(e.read().decode("utf-8"))
                if "message" in error_body:
                    msg = error_body["message"]
            except Exception:
                pass
            messagebox.showerror("Server Error", f"Code {e.code}: {msg}")
            return "ERROR"
        except URLError:
            messagebox.showerror("Connection Error", f"Could not reach server at {BASE_URL}.\nEnsure your backend daemon is running.")
            return "ERROR"

    def load_watchlist(self, status_filter=None):
        """Fetches items and populates local database management panel."""
        endpoint = "/movies"
        if status_filter:
            endpoint += f"?status={status_filter}"
            
        result = self.make_network_request(endpoint)
        if result == "ERROR":
            return
            
        # Clear existing table content
        for item in self.watchlist_tree.get_children():
            self.watchlist_tree.delete(item)
            
        if not result:
            return
            
        for m in result:
            # Handle standard formatting logic with explicit fallback values for nulls
            m_id = m.get("id")
            title = m.get("title", "")
            year = m.get("release_year") if m.get("release_year") is not None else "N/A"
            status = m.get("status", "pending")
            
            raw_genres = m.get("genre", [])
            if isinstance(raw_genres, str):
                try:
                    raw_genres = json.loads(raw_genres)
                except json.JSONDecodeError:
                    pass
            genre_str = ", ".join(raw_genres) if isinstance(raw_genres, list) else str(raw_genres)
            
            self.watchlist_tree.insert("", tk.END, values=(m_id, title, year, genre_str, status))

    def search_tmdb(self):
        """Proxies string request out to external data layer."""
        query = self.search_entry.get().strip()
        if not query:
            return
            
        # Clear existing search results table
        for item in self.search_tree.get_children():
            self.search_tree.delete(item)
            
        result = self.make_network_request("/movies/search", method="POST", data={"query": query})
        if result == "ERROR" or not result:
            return
            
        for r in result:
            tmdb_id = r.get("tmdb_id")
            title = r.get("title", "")
            year = r.get("release_year") if r.get("release_year") is not None else "N/A"
            genres = ", ".join(r.get("genre", [])) if isinstance(r.get("genre"), list) else ""
            
            # Embed real array structure inside row attachment reference data tags
            self.search_tree.insert("", tk.END, values=(tmdb_id, title, year, genres), tags=(json.dumps(r.get("genre", [])),))

    def add_selected_movie(self):
        """Pushes structured item information downstream into storage layer."""
        selected = self.search_tree.selection()
        if not selected:
            messagebox.showwarning("Selection Missing", "Please select a movie from the search results table first.")
            return
            
        values = self.search_tree.item(selected[0], "values")
        tags = self.search_tree.item(selected[0], "tags")
        
        tmdb_id = int(values[0])
        title = values[1]
        year = int(values[2]) if values[2] != "N/A" else None
        genres = json.loads(tags[0]) if tags else []
        
        payload = {
            "tmdb_id": tmdb_id,
            "title": title,
            "release_year": year,
            "genre": genres
        }
        
        result = self.make_network_request("/movies", method="POST", data=payload)
        if result != "ERROR":
            # Automatically refresh active database visual lists
            self.load_watchlist()

    def toggle_movie_status(self):
        """Flips processing values tracking state flags between pending and watched."""
        selected = self.watchlist_tree.selection()
        if not selected:
            messagebox.showwarning("Selection Missing", "Please select a movie from your active watchlist first.")
            return
            
        values = self.watchlist_tree.item(selected[0], "values")
        movie_id = values[0]
        current_status = values[4]
        
        # Toggle state logic path
        new_status = "watched" if current_status == "pending" else "pending"
        
        result = self.make_network_request(f"/movies/{movie_id}", method="PATCH", data={"status": new_status})
        if result != "ERROR":
            self.load_watchlist()

if __name__ == "__main__":
    app = WatchlistApp()
    app.mainloop()
