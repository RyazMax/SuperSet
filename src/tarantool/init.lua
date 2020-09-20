local config = require('config').config
local ddl = require('ddl')

local user = require('src.tarantool.user')
local project = require('src.tarantool.project')
local session = require('src.tarantool.session')
local schema = require('src.tarantool.schema')
local grant = require('src.tarantool.grant')
local labeledtask = require('src.tarantool.labeledtask')
local tasks = require('src.tarantool.tasks')

box.cfg{
    listen = config.listen_port,
    work_dir = config.work_dir,
}

local function init()
    local ok, err
    ok, err = ddl.set_schema(schema.schema)
    if err then
        error(err)
    end
    ok, err = ddl.set_schema(user.schema)
    if err then
        error(err)
    end
    ok, err = ddl.set_schema(session.schema)
    if err then
        error(err)
    end

    box.schema.user.create('go',  { password = 'go', if_not_exists = true })
    box.schema.user.grant('go', 'read,write', 'universe', nil, { if_not_exists = true })

    grant.init()
    labeledtask.init()
    tasks.init()
    project.init()
end

init()
require('console').start()