local config = require('config').config
local ddl = require('ddl')

local user = require('src.tarantool.user')
local project = require('src.tarantool.project')

box.cfg{
    listen = config.listen_port,
    work_dir = config.work_dir,
}

local function init()
    ddl.set_schema(user.schema)
    ddl.set_schema(project.schema)

    box.schema.user.create('go',  { password = 'go', if_not_exists = true })
    box.schema.user.grant('go', 'read,write', 'universe', nil, { if_not_exists = true })
end

init()
require('console').start()