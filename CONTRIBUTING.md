# Contributing to KubeAuth

Thank you for your interest in contributing to KubeAuth!

## How to contribute

1. Fork the repository
2. Create a new branch: `git checkout -b feat/your-feature`
3. Make your changes
4. Run the build: `go build ./...`
5. Commit: `git commit -m "feat: your feature description"`
6. Push: `git push origin feat/your-feature`
7. Open a Pull Request

## Commit convention

We follow conventional commits:
- `feat:` new feature
- `fix:` bug fix
- `docs:` documentation
- `ci:` CI/CD changes
- `chore:` maintenance

## Development setup

```bash
git clone https://github.com/Abderraoufzekkour/KubeAuth.git
cd KubeAuth
go mod tidy
go build -o bin/kubeauth ./cmd/kubeauth
./bin/kubeauth --help
```

## Questions

Open an issue on GitHub.