module ai-coding-training

go 1.24.0

require (
	github.com/golang-migrate/migrate/v4 v4.19.1 // indirect
	github.com/lib/pq v1.12.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	raft-stash v0.0.0
)

replace raft-stash => ./raft-stash/src
