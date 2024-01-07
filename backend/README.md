# Pulsechat Microservices Architecture

Pulsechat is a microservices-based messaging app being developed in Golang, designed to provide secure and efficient end-to-end encrypted messaging, group messaging, and file sharing.

## Services

- **auth:** Handles user authentication with various providers (Google, Facebook, GitHub) and conventional email-password authentication.
- **messaging:** Manages the core messaging functionality, including sending and receiving messages.
- **notification:** Provides notification services for various events within the app.
- **storage:** Responsible for managing file storage and sharing functionalities.
- **user:** Manages user-related operations, including profile management.
