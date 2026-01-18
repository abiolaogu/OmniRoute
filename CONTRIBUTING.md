# OmniRoute - Contributing Guide

## Welcome Contributors! 🎉

Thank you for your interest in contributing to OmniRoute. This guide will help you get started.

---

## Code of Conduct

We are committed to providing a welcoming and inclusive environment. Please:

- Be respectful and considerate
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Accept differing viewpoints gracefully

---

## Getting Started

### 1. Fork the Repository

```bash
# Fork via GitHub UI, then clone your fork
git clone https://github.com/YOUR-USERNAME/omniroute.git
cd omniroute

# Add upstream remote
git remote add upstream https://github.com/omniroute/omniroute.git
```

### 2. Set Up Development Environment

```bash
# Install dependencies
make setup

# Start infrastructure
docker-compose up -d postgres redis kafka

# Verify setup
make test
```

### 3. Create a Branch

```bash
# Get latest main
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-feature-name
```

---

## Development Workflow

### 1. Follow TDD (Test-Driven Development)

```go
// 1. Write a failing test first
func TestNewFeature_WhenCondition_ShouldDoSomething(t *testing.T) {
    // Arrange
    sut := NewSystemUnderTest()
    
    // Act
    result := sut.DoSomething()
    
    // Assert
    assert.Equal(t, expected, result)
}

// 2. Write minimal code to pass
// 3. Refactor while keeping tests green
```

### 2. Follow DDD Principles

```go
// Domain models encapsulate business logic
type Order struct {
    // ... fields
}

// Rich domain methods
func (o *Order) Cancel(reason string) error {
    if !o.CanBeCancelled() {
        return ErrOrderCannotBeCancelled
    }
    o.Status = OrderStatusCancelled
    o.CancelReason = reason
    o.AddEvent(OrderCancelledEvent{...})
    return nil
}
```

### 3. Run Quality Checks

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Run all checks
make check
```

---

## Pull Request Process

### 1. Prepare Your PR

```bash
# Ensure tests pass
make test

# Ensure lint passes
make lint

# Update your branch
git fetch upstream
git rebase upstream/main
```

### 2. Write a Good PR Description

```markdown
## Summary
Brief description of what this PR does.

## Changes
- Added X feature to Y service
- Updated Z model to include A

## Testing
- Added unit tests for X
- Manually tested Y workflow

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Lint passes
- [ ] Ready for review
```

### 3. Address Review Feedback

- Respond to all comments
- Push fixes as new commits
- Request re-review when ready

---

## Coding Standards

### Go Code Style

```go
// ✅ Good: Clear, descriptive names
func (r *OrderRepository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*Order, error) {
    // ...
}

// ❌ Bad: Unclear abbreviations
func (r *OrderRepo) FndByCust(ctx context.Context, cid uuid.UUID) ([]*Order, error) {
    // ...
}
```

### Error Handling

```go
// ✅ Good: Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("find order by customer %s: %w", customerID, err)
}

// ❌ Bad: Raw error return
if err != nil {
    return nil, err
}
```

### Logging

```go
// ✅ Good: Structured logging with context
logger.Info("order created",
    zap.String("order_id", order.ID.String()),
    zap.String("customer_id", order.CustomerID.String()),
    zap.Decimal("total", order.Total),
)

// ❌ Bad: Printf style
log.Printf("order %s created for customer %s", orderID, customerID)
```

---

## Types of Contributions

### 🐛 Bug Fixes

1. Check existing issues first
2. Create issue if none exists
3. Reference issue in PR

### ✨ Features

1. Discuss in issue first
2. Get approval from maintainers
3. Implement with tests

### 📚 Documentation

1. Fix typos, improve clarity
2. Add examples
3. Update outdated content

### 🧪 Tests

1. Add missing test coverage
2. Improve test reliability
3. Add integration tests

---

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation |
| `style` | Formatting (no code change) |
| `refactor` | Code restructuring |
| `test` | Adding tests |
| `chore` | Maintenance |
| `perf` | Performance improvement |

### Examples

```
feat(pricing): add volume discount calculation

Implements tiered pricing based on order quantity.
Supports up to 5 tiers with configurable thresholds.

Closes #123
```

```
fix(payment): handle timeout in bank gateway

- Add 30s timeout for bank API calls
- Implement retry with exponential backoff
- Add circuit breaker for repeated failures

Fixes #456
```

---

## Issue Guidelines

### Bug Report Template

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
1. Go to '...'
2. Click on '...'
3. See error

**Expected behavior**
What you expected to happen.

**Environment**
- OS: [e.g., macOS 14.2]
- Browser: [if applicable]
- Version: [e.g., v1.2.3]

**Additional context**
Any other relevant information.
```

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
A clear description of the problem.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Any alternative solutions or features.

**Additional context**
Any other context or screenshots.
```

---

## Recognition

Contributors are recognized in:

- CONTRIBUTORS.md file
- Release notes
- Documentation acknowledgments

---

## Getting Help

- 💬 [Discord Community](https://discord.gg/omniroute)
- 📧 [Email](mailto:dev@omniroute.io)
- 📖 [Documentation](/docs)
- 🐛 [Issues](https://github.com/omniroute/omniroute/issues)

---

Thank you for contributing! 🙏
