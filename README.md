# yantra

A factory that makes software factories.

A spec and a public test suite go in.
A production line comes out: failing tests become tickets, coding agents make them pass, every pull request gets a Qodo review, and a gate merges or sends the change to a human.

Built for the Qodo Software Factory Hackathon, November 2026.

## Status

Seed commit.
Nothing is built yet.
Every feature after this commit arrives as a GitHub issue and leaves as a reviewed pull request.

## Run

```
go build ./cmd/yantra
./yantra version
```

## Verify

```
go vet ./...
go test -race ./...
```

## License

MIT
