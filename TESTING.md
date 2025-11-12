# Testing Documentation

## Test Coverage

This project maintains high test coverage for critical business logic components.

### Coverage Goals

- **Core packages (config, storage, timeparse)**: >80% coverage ✅
- **API clients (jira, tempo)**: Structure and unit tests ✅
- **Commands (cmd)**: Integration tests not included (interactive CLI)
- **UI (internal/ui)**: Not tested (uses third-party survey library)

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run specific package tests
go test -v ./internal/timeparse/
go test -v ./internal/storage/
go test -v ./internal/config/

# Generate coverage HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Coverage by Package

| Package   | Coverage | Notes                                     |
| --------- | -------- | ----------------------------------------- |
| config    | 83.0%    | Configuration loading and validation      |
| storage   | 81.8%    | SQLite database operations                |
| timeparse | 89.2%    | Time format parsing and rounding          |
| jira      | 8.2%     | Structure tests (HTTP calls not mocked)   |
| tempo     | 13.8%    | Structure tests (HTTP calls not mocked)   |
| ui        | 0.0%     | Interactive prompts (integration testing) |
| cmd       | 0.0%     | CLI commands (integration testing)        |

### What's Tested

#### Config Package ✅
- Configuration validation (all required fields)
- Shortcut lookup
- Label filtering
- YAML parsing
- File loading edge cases

#### Storage Package ✅
- Database initialization
- Time entry CRUD operations
- Today's entries retrieval
- Unsynced entries tracking
- Total seconds calculation
- SQLite schema creation

#### Timeparse Package ✅
- Multiple time format parsing (2h 30m, 2.5h, 150m)
- Rounding to nearest 5 minutes
- Time validation
- Format conversion
- Edge cases (zero, negative, invalid)

#### Jira Package ✅
- Client initialization
- Structure definitions
- Timeout configuration
- URL handling

#### Tempo Package ✅
- Client initialization
- Request/response structures
- Format utilities

### What's Not Tested

1. **HTTP API Calls**: Jira and Tempo API clients make real HTTP calls. These would require:
   - HTTP mocking libraries (httptest)
   - Mock servers
   - Response fixtures
   
2. **Interactive UI**: The survey-based prompts require user interaction:
   - Task selection
   - Time entry prompts
   - Label selection
   - Confirmation dialogs

3. **CLI Commands**: The cobra commands require integration testing:
   - `tasklog log`
   - `tasklog sync`
   - `tasklog summary`
   - `tasklog init`

4. **Main Package**: Entry point with minimal logic

### Future Testing Improvements

1. **API Mocking**: Add httptest-based tests for Jira/Tempo clients
2. **Integration Tests**: Add end-to-end testing for CLI commands
3. **Benchmark Tests**: Performance testing for database operations
4. **Fuzzing**: Fuzz testing for time parsing

### Manual Testing Checklist

For interactive features, use this manual testing checklist:

- [ ] `tasklog init` creates config file
- [ ] `tasklog log` shows in-progress tasks
- [ ] Task search works correctly
- [ ] Time formats are accepted and rounded
- [ ] Labels are filtered correctly
- [ ] Shortcuts work as expected
- [ ] Summary shows today's entries
- [ ] Sync retries failed entries
- [ ] Database is created on first run
- [ ] Config validation errors are clear

## Test Philosophy

We focus test coverage on:
1. **Business logic** (parsing, validation, calculations)
2. **Data persistence** (database operations)
3. **Configuration** (loading, validation)

We accept lower coverage for:
1. **External integrations** (HTTP APIs)
2. **User interaction** (CLI prompts)
3. **Infrastructure** (file I/O, network calls)

This approach ensures reliability of core functionality while keeping tests maintainable and fast.
