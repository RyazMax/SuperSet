start_tarantool:
	tarantool src/tarantool/init.lua

clean:
	rm src/tarantool/data/* -rf

test:
	go test ./...