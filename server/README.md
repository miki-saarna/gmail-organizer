connect to GMAIL API

GMAIL API endpoints

find all emails:

- filter by sender email address

deletion/archive:

- filter by sender email address
- filter by subject line

unsubscribe:

- by email address?

later:

- incorporating AI

Working CLI:
find all emails:

- Input email address as first argument. Example: go run ./cmd/main.go "helloworld@gmail.com"

DIRECTIONS

- When changing scope, MUST delete `token.json` file so new one can be recreated
- Copy `deletionList.example.json` and rename `deletionList.json`. Fill with relevant email address values
