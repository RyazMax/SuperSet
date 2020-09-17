start_tarantool:
	mkdir -p src/tarantool/data
	tarantool src/tarantool/init.lua

clean:
	rm src/tarantool/data/* -rf
	rm src/static/media/* -rf

test:
	go test ./...