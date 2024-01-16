
.PHONY: mock rebuild_db run
# required tools:
# sqlite3
# mkcert https://github.com/FiloSottile/mkcert && mkcert -install && mkcert localhost
# go
# bun

database: database/database.db

database/database.db: database/database.sql
	@echo "Creating database..."
	rm -f database/database.db
	sqlite3 database/database.db < database/database.sql
	@echo "Database created."

rebuild_db: 
	@echo "Rebuilding database..."
	rm -f database/database.db
	sqlite3 database/database.db < database/database.sql
	@echo "Database rebuilt."

mock: mock_users mock_requests mock_friends mock_messages
	@echo "Mock data created."

mock_users: 
	@echo "Creating mock users..."
	sqlite3 database/database.db < database/mock/users.sql 
	@echo "Mock users created."

mock_requests:
	@echo "Creating mock requests..."
	sqlite3 database/database.db < database/mock/requests.sql
	@echo "Mock requests created."

mock_friends:
	@echo "Creating mock friends..."
	sqlite3 database/database.db < database/mock/friends.sql
	@echo "Mock friends created."

mock_messages:
	@echo "Creating mock messages..."
	sqlite3 database/database.db < database/mock/messages.sql
	@echo "Mock messages created."

create_cert:
	openssl req -x509 -newkey rsa:4096 -keyout localhost.key -out localhost.crt -days 365 -nodes -subj "/CN=localhost"

build: 
	@echo "Building server..."
	cd server && CGO_ENABLED=1 go build -o server server.go
	@echo "Building React app..."
	cd app && bun run build
	@echo "Build complete."

start:
	@echo "Starting server..."
	cd server && ./server
	@echo "Server started."


run: 
	@echo "Running server..."
	cd server && CGO_ENABLED=1 go run server.go