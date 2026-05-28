# Movies To Watch API Specification

## Endpoints

### 1. Fetch Watchlist
Retrieves saved movies from the local database tracker, optionally filtered by status and/or genres.

* **URL:** `/movies`
* **Method:** `GET`
* **Query Parameters:** 
    * `status` (optional): `pending` or `watched`
    * `genres` (optional): A comma-separated list of case-sensitive genre strings (e.g., `?genres=Thriller,Mystery`). When provided. The server returns movies that contain any of the specified genres.
* **Response Codes:**
  * `200 OK`: Request successful.
* **Response Body (Received):**
```json
[
  {
    "id": 1,
    "tmdb_id": 27205,
    "title": "Inception",
    "release_year": 2010,
    "genre": ["Action", "Science Fiction", "Adventure"],
    "status": "pending"
  }
]
```

### **2. Search Remote Catalog**

Proxies a movie title search string to the third-party TMDB catalog API. Does not alter the database.

* **URL:** /movies/search  
* **Method:** POST  
* **Response Codes:**  
  * 200 OK: Search successfully processed.  
* **Request Body (Sent):**

```json
{  
  "query": "Interstellar"  
}
```

* **Response Body (Received):**

```json
[  
  {  
    "tmdb_id": 157336,  
    "title": "Interstellar",  
    "release_year": 2014,  
    "genre": ["Adventure", "Drama", "Science Fiction"]  
  }  
]
```

### **3. Add Movie to Watchlist**

Inserts a discovered movie into the local persistent tracking database.

* **URL:** /movies  
* **Method:** POST  
* **Response Codes:**  
  * 201 Created: Movie successfully stored. Returns an empty body.  
  * 409 Conflict: The tmdb_id already exists within the unique database index.  
* **Request Body (Sent):**

```json  
{  
  "tmdb_id": 157336,  
  "title": "Interstellar",  
  "release_year": 2014,  
  "genre": ["Adventure", "Drama", "Science Fiction"]  
}
```

* **Response Body (Received):**  
  *None (Empty Response)*

### **4. Update Watchlist Status**

Modifies the tracking state flag of an existing database entry using its internal sequencing key.

* **URL:** /movies/{id}  
* **Method:** PATCH  
* **URL URL Parameters:**  
  * {id}: Internal database auto-increment integer ID  
* **Response Codes:**  
  * 200 OK: Status successfully toggled. Returns an empty body.  
  * 404 Not Found: The requested internal tracking {id} does not exist.  
* **Request Body (Sent):**

```json  
{  
  "status": "watched"  
}
```

* **Response Body (Received):**  
  *None (Empty Response)*

### **5. Remove A Movie From Watchlist**

Removes a movie from from the database.

* **URL:** /movies/{id}
* **Method:** DELETE  
* **Response Codes:**  
  * 200 OK: Movie successfully deleted. Returns an empty body.  
  * 404 Not Found: The requested internal tracking {id} does not exist.  
* **Request Body (Sent):**

  *None (Empty Response)*

* **Response Body (Received):**  
  *None (Empty Response)*

## **Global Handshake Validation (CORS Preflight)**

Browsers automatically prepend restricted cross-origin actions with an environment handshake verification route.

* **URL:** *All Endpoints*  
* **Method:** OPTIONS  
* **Response Codes:**  
  * 204 No Content or 200 OK  
* **Required Headers (Returned by Server):**

```text
Access-Control-Allow-Origin: *  
Access-Control-Allow-Methods: GET, POST, PATCH, DELETE, OPTIONS  
Access-Control-Allow-Headers: Content-Type
```

* **Response Body (Received):**  
  *None (Empty Response)*
