# Gator

Gator is a guided RSS project from Boot.Dev. It allows users to manage RSS feeds, follow feeds, and browse posts using a command-line interface.

## Prerequisites

To use Gator, you need:

- **PostgreSQL**: Ensure you have a running PostgreSQL instance.
- **Go**: Install Go (version 1.24.1 or later).

## Installation

To install Gator, run the following command in your terminal:

```bash
go install github.com/GLobyNew/gator@latest
```

## Configuration

Gator requires a configuration file named `.gatorconfig.json` in your home directory. The file should include the following fields:

```json
{
  "db_url": "your_postgres_connection_string",
  "current_user_name": "your_username"
}
```

Replace `your_postgres_connection_string` with your PostgreSQL connection string and `your_username` with your desired username.

## Usage

Run Gator with the following commands:

### User Management

- **Register a new user**:
  ```bash
  gator register <username>
  ```
- **Login as an existing user**:
  ```bash
  gator login <username>
  ```
- **List all users**:
  ```bash
  gator users
  ```
- **Reset the database**:
  ```bash
  gator reset
  ```

### Feed Management

- **Add a new feed**:
  ```bash
  gator addfeed <feed_name> <feed_url>
  ```
- **List all feeds**:
  ```bash
  gator feeds
  ```
- **Follow a feed**:
  ```bash
  gator follow <feed_url>
  ```
- **List followed feeds**:
  ```bash
  gator following
  ```
- **Unfollow a feed**:
  ```bash
  gator unfollow <feed_url>
  ```

### Browsing Posts

- **Browse posts**:
  ```bash
  gator browse [limit]
  ```
  The `limit` parameter is optional and defaults to 2.

### Aggregation

- **Start feed aggregation**:
  ```bash
  gator agg <time_between_requests>
  ```
  Replace `<time_between_requests>` with a duration (e.g., `1m` for 1 minute).

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE).

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## Acknowledgments

This project is part of the Boot.Dev guided learning program.