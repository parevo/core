# Contributing

Parevo Core is an open source project. Contributions are welcome.

## Issue-first workflow

1. **Open an issue first**: Before proposing new features or major changes.
2. **Discuss**: Align on design with maintainers and the community.
3. **Submit a PR**: Share your implementation via pull request after approval.

Bug fixes and documentation improvements can be submitted directly as PRs.

## Issue types

- **Bug Report**: Unexpected behavior or errors
- **Feature Request**: New feature proposal
- **Documentation**: Documentation improvements
- **Question**: Questions or discussion (prefer Discussions)

## Pull request process

1. Fork and create a feature branch (`git checkout -b feature/amazing-feature`)
2. Make your changes
3. Run tests: `go test ./...`
4. Run lint: `make lint` or `golangci-lint run`
5. Use meaningful commit messages
6. Open a PR and reference the related issue

## Code standards

- Must pass `go fmt` and `go vet`
- Add documentation for new public APIs
- Consider adding an example under `examples/` for new features

## Adding examples

New feature examples live under `examples/` in separate folders:

```
examples/
  my-feature/
    main.go
```

Add a short description to `examples/README.md`.
