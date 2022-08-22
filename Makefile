migrate-up:
	@migrate -path ./migrations -database "mysql://root:secret@tcp(localhost:3306)/cake-store" -verbose up
migrate-down:
	@migrate -path ./migrations -database "mysql://root:secret@tcp(localhost:3306)/cake-store" -verbose down