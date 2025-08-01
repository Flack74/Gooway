# 🔗 Gooway

![Go Version](https://img.shields.io/badge/Go-1.24.5-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white)
![YAML](https://img.shields.io/badge/YAML-CB171E?style=for-the-badge&logo=yaml&logoColor=white)
![JSON](https://img.shields.io/badge/JSON-000000?style=for-the-badge&logo=json&logoColor=white)

> 🚀 A powerful and flexible URL shortener service built with Go, supporting multiple data sources including SQLite database, YAML, and JSON configurations.

## ✨ Features

- 🗄️ **Multiple Data Sources**: Support for SQLite database, YAML, and JSON files
- 📊 **Rich Database**: Pre-loaded SQLite database with 500+ real website URLs
- 🔄 **Flexible Fallback System**: Hierarchical URL resolution with graceful fallbacks
- ⚡ **High Performance**: Built with Go's native HTTP server for optimal speed
- 🛠️ **Easy Configuration**: Simple command-line flags for different data sources
- 📁 **Modular Design**: Clean separation of concerns with reusable handlers
- 🔒 **Production Ready**: Robust error handling and database connection management

## 🏗️ Architecture

```
url_shortner/
├── 📁 main/
│   └── main.go           # Application entry point
├── 📁 urlshort/
│   └── handler.go        # URL handling logic
├── 📁 data/
│   ├── urls.db          # SQLite database with 500+ real websites
│   ├── urls.yaml        # YAML configuration (24 sample URLs)
│   └── urls.json        # JSON configuration (40 sample URLs)
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
└── README.md            # This file
```

## 🚀 Quick Start

### Prerequisites

- 🐹 **Go 1.24.5+** installed on your system
- 📊 **SQLite3** support (automatically handled by Go driver)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Flack74/Gooway.git
   cd Gooway
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the application:**
   ```bash
   go run main/main.go
   ```

The server will start on `http://localhost:8080` 🌐

## 💻 Usage

### Basic Usage

Start the server with default SQLite database:
```bash
go run main/main.go
```

### Using YAML Configuration

```bash
go run main/main.go -yaml=data/urls.yaml
```

### Using JSON Configuration

```bash
go run main/main.go -json=data/urls.json
```

### Using Both YAML and JSON

```bash
go run main/main.go -yaml=data/urls.yaml -json=data/urls.json
```

## 🎯 How It Works

The URL shortener uses a **hierarchical fallback system**:

1. 🥇 **JSON Handler** (if specified) - Highest priority
2. 🥈 **YAML Handler** (if specified) - Second priority
3. 🥉 **SQLite Handler** - Default fallback
4. 🔄 **Default Mux** - Final fallback (returns "Hello, world!")

### Example URL Mappings

| Short URL | Redirects To | Source |
|-----------|--------------|--------|
| `localhost:8080/google` | `https://www.google.com` | SQLite DB |
| `localhost:8080/github` | `https://github.com` | All sources |
| `localhost:8080/spotify` | `https://spotify.com` | SQLite DB |
| `localhost:8080/netflix` | `https://www.netflix.com` | All sources |
| `localhost:8080/amazon` | `https://amazon.com` | SQLite DB |
| `localhost:8080/microsoft` | `https://microsoft.com` | SQLite DB |

*The SQLite database contains 500+ real website URLs across multiple categories!*

## 📄 Configuration Files

### YAML Format (`data/urls.yaml`)

```yaml
- path: /google
  url: https://www.google.com
- path: /github
  url: https://github.com
- path: /golang
  url: https://golang.org
```

### JSON Format (`data/urls.json`)

```json
[
  { "path": "/google", "url": "https://www.google.com" },
  { "path": "/github", "url": "https://github.com" },
  { "path": "/golang", "url": "https://golang.org" }
]
```

### SQLite Database (`data/urls.db`)

The SQLite database comes **pre-loaded with 500+ real website URLs** covering:

🌐 **Categories Included:**
- **Technology**: Google, Microsoft, Apple, GitHub, Adobe, Tesla
- **Social Media**: Facebook, Instagram, Twitter/X, LinkedIn, TikTok, Discord
- **E-commerce**: Amazon, eBay, Shopify, Alibaba, Flipkart
- **Streaming**: Netflix, Spotify, YouTube, Disney+, Hulu, Twitch
- **Banking**: Chase, PayPal, Visa, Mastercard, Bank of America
- **News**: BBC, CNN, Reuters, New York Times, Forbes
- **Travel**: Booking.com, Airbnb, Expedia, TripAdvisor
- **Food Delivery**: Uber Eats, DoorDash, Zomato, Swiggy
- **Education**: Universities and online learning platforms
- **Government**: Official government websites from various countries

**Database Structure:**
```sql
CREATE TABLE urls (
    path TEXT NOT NULL PRIMARY KEY,
    url  TEXT NOT NULL
);
```

**Sample Database Entries:**
```sql
INSERT INTO urls (path, url) VALUES ('/google', 'https://google.com');
INSERT INTO urls (path, url) VALUES ('/amazon', 'https://amazon.com');
INSERT INTO urls (path, url) VALUES ('/spotify', 'https://spotify.com');
INSERT INTO urls (path, url) VALUES ('/netflix', 'https://netflix.com');
INSERT INTO urls (path, url) VALUES ('/github', 'https://github.com');
-- ... and 495 more real websites!
```

✨ **Key Features of the Database:**
- 🏢 **Real websites only** - No fictional entries
- 🧹 **Clean paths** - Simple, readable format (e.g., `/google`, `/spotify`)
- 🚫 **No duplicates** - Each path is unique
- ⚡ **Production-ready** - All websites are active and well-known
- 🌍 **Global coverage** - Websites from multiple countries and industries

## 🔧 API Reference

### Core Handlers

#### `MapHandler(pathsToUrls map[string]string, fallback http.Handler)`
- **Purpose**: Creates an HTTP handler from a map of paths to URLs
- **Parameters**:
  - `pathsToUrls`: Map containing path-to-URL mappings
  - `fallback`: Handler to call when path is not found
- **Returns**: `http.HandlerFunc`

#### `YAMLHandler(yml []byte, fallback http.Handler)`
- **Purpose**: Creates a handler from YAML configuration
- **Parameters**:
  - `yml`: YAML bytes containing URL mappings
  - `fallback`: Fallback handler
- **Returns**: `http.HandlerFunc, error`

#### `JSONHandler(jsondata []byte, fallback http.Handler)`
- **Purpose**: Creates a handler from JSON configuration
- **Parameters**:
  - `jsondata`: JSON bytes containing URL mappings
  - `fallback`: Fallback handler
- **Returns**: `http.HandlerFunc, error`

#### `SqliteHandler(db *sql.DB, fallback http.Handler)`
- **Purpose**: Creates a handler from SQLite database
- **Parameters**:
  - `db`: SQLite database connection
  - `fallback`: Fallback handler
- **Returns**: `http.HandlerFunc, error`

## 🛠️ Development

### Building

```bash
go build -o url-shortener main/main.go
```

### Running Tests

```bash
go test ./...
```

### Code Structure

The project follows Go best practices with clear separation of concerns:

- **`main/main.go`**: Application bootstrap and configuration
- **`urlshort/handler.go`**: Core URL handling logic
- **`data/`**: Configuration and database files

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. 🍴 **Fork** the repository
2. 🌱 **Create** a feature branch (`git checkout -b feature/AmazingFeature`)
3. 💾 **Commit** your changes (`git commit -m 'Add some AmazingFeature'`)
4. 📤 **Push** to the branch (`git push origin feature/AmazingFeature`)
5. 🔄 **Open** a Pull Request

### Development Guidelines

- ✅ Follow Go formatting standards (`go fmt`)
- 📝 Add tests for new functionality
- 📚 Update documentation as needed
- 🧹 Keep code clean and well-commented

## 📈 Performance

- **Response Time**: < 1ms for cached URLs
- **Throughput**: Supports thousands of concurrent requests
- **Memory Usage**: Minimal footprint with efficient map-based lookups
- **Database**: SQLite provides excellent performance for small to medium datasets

## 🔒 Security Considerations

- ✅ Input validation for URL paths
- 🛡️ SQL injection protection through prepared statements
- 🔄 Proper HTTP status codes (302 for redirects)
- 🚫 No sensitive data exposure in error messages

## 🌟 Roadmap

- [ ] 📊 Admin dashboard for URL management
- [ ] 📈 Click tracking and analytics
- [ ] 🔐 Authentication and authorization
- [ ] 🏷️ Custom short URL aliases
- [ ] ⏰ URL expiration dates
- [ ] 🌐 REST API endpoints
- [ ] 📱 Rate limiting
- [ ] 🎨 Web UI for URL creation

## 📚 Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `gopkg.in/yaml.v3` | `v3.0.1` | YAML parsing |
| `github.com/mattn/go-sqlite3` | `v1.14.30` | SQLite driver |

## 🐛 Troubleshooting

### Common Issues

**Server won't start:**
- ✅ Check if port 8080 is available
- 🔍 Verify Go installation
- 📁 Ensure proper file permissions

**Database errors:**
- 🗄️ Check if `data/urls.db` exists and is readable
- 🔧 Verify SQLite table structure
- 🚫 Ensure no database locks

**YAML/JSON parsing errors:**
- 📝 Validate file syntax
- 🔍 Check file paths and permissions
- 📋 Verify data structure matches expected format

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- 🐹 **Gophercises** - Original exercise inspiration
- 🏗️ **Go Community** - Amazing ecosystem and support
## 📞 Support

- 📧 **Email**: [puspendrachawlax@gmail.com]
- 🐛 **Issues**: [GitHub Issues](https://github.com/Flack74/Gooway/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/Flack74/Gooway/discussions)

---

<div align="center">

**Made with ❤️ by Flack**

</div>
