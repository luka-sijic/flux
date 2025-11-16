# Flux Messaging App
## Apps
Go Websocket server\
Go HTTP server\
NextJS frontend

# Features
Paired with go-auth-v2 for user registration and authentication
- websocket w/ redis streams impl (kafka and scylladb were my goal but fell short)
- bloom filters for early auth rejection
- sharded postgres using snowflake IDs
- friend system with a multi db setup
- user profiles w/ pfps (s3)
- JWT cookie system

## Databases
Postgres\
Redis
#### TBD:
ScyallDB\
Kafka

