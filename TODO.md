# TODO: Gin Framework Features Checklist

## Core Features
- [x] **Routing**: Define and handle HTTP routes (GET, POST, PUT, DELETE, etc.).
- [x] **Middleware**: Implement custom middleware for request processing.
- [ ] **Logging**: Customize logging format and output (e.g., JSON logging).
  - global logger or inject logger
- [x] **Error Handling**: Centralized error handling for consistent responses.
- [ ] **Container**: Dependency injection for modular applications. Where to provide? [fx](https://uber-go.github.io/fx/index.html)
 
## Advanced Features
- [ ] **Authentication**: Implement JWT, OAuth, or Basic Authentication.
- [ ] **Validation**: Validate request payloads using struct tags or custom validators.
- [ ] **Database Integration**: Integrate with databases like MySQL, PostgreSQL, or MongoDB.
- [ ] **Caching**: Implement caching mechanisms (e.g., Redis) for improved performance.
- [ ] **wire/dig inject**: [wire](https://github.com/google/wire) 

## Deployment & Maintenance
- [ ] **Dockerization**: Create Dockerfiles for containerized deployment.
- [ ] **Health Checks**: Add health check endpoints for monitoring.
- [ ] **Metrics**: Integrate with Prometheus or other monitoring tools.
- [ ] **CI/CD**: Set up CI/CD pipelines for automated testing and deployment.

## Testing
- [ ] **Unit Tests**: Write unit tests for individual components.
- [ ] **Integration Tests**: Write integration tests for API endpoints.
- [ ] **Mocking**: Use mocking libraries for testing external dependencies.

## Documentation
- [x] **README**: Maintain an up-to-date README with setup instructions.
- [ ] **API Docs**: Generate and maintain API documentation.
- [x] **Instructions for gin bind Usage**：[doc/bind.md](docs/bind.md)
- [x] **bind struct tag**：[doc/bind_tag.md](docs/bind_tag.md)


