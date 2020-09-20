local ddl = require('ddl')

local project = {}

project.schema = {
    spaces = {
        projects = {
            engine = 'memtx',
            is_local = false,
            temporary = false,
            format = {
                {name = 'ID', is_nullable = false, type = 'unsigned'},
                {name = 'Name', is_nullable = false, type = 'string'},
                {name = 'OwnerID', is_nullable = false, type = 'unsigned'},
            },
            indexes = {{
                name = 'ID',
                type = 'HASH',
                unique = true,
                parts = {
                    {path = 'ID', is_nullable = false, type = 'unsigned'}
                }
            }, {
                name = 'name',
                type = 'TREE',
                unique = true,
                parts = {
                    {path = 'Name', is_nullable = false, type = 'string'}
                }
            }, {
                name = 'ownerid',
                type = 'TREE',
                unique = false,
                parts = {
                    {path = 'OwnerID', is_nullable = false, type = 'unsigned'}
                }
            }},
        },
    },
}

function get_allowed_projects(id)
    local res = {}
    for _, t in box.space.project_grants.index.UserID:pairs(id) do
        table.insert(res, box.space.projects:get(t["ProjectID"]))
    end
    return res
end

function project.init()
    local ok, err = ddl.set_schema(project.schema)
    if err ~= nil then
        error(err)
    end
end

return project