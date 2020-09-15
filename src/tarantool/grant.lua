local ddl = require('ddl')

local grant = {}

grant.schema = {
    spaces = {
        project_grants = {
            engine = 'memtx',
            is_local = false,
            temporary = false,
            format = {
                {name = 'ProjectID', is_nullable = false, type = 'unsigned'},
                {name = 'UserID', is_nullable = false, type = 'unsigned'},
            },
            indexes = {{
                name = 'ID',
                type = 'HASH',
                unique = true,
                parts = {
                    {path = 'ProjectID', is_nullable = false, type = 'unsigned'},
                    {path = 'UserID', is_nullable = false, type = 'unsigned'}
                }
            },{
                name = 'ProjectID',
                type = 'TREE',
                unique = false,
                parts = {
                    {path = 'ProjectID', is_nullable = false, type = 'unsigned'},
                }
            },{
                name = 'UserID',
                type = 'TREE',
                unique = false,
                parts = {
                    {path = 'UserID', is_nullable = false, type = 'unsigned'},
                }
            }
        },
        },
    },
}

function delete_grants_by_project_id(id)
    for _,t in box.space.project_grants.index.ProjectID:pairs() do
        box.space.project_grants:delete(t)
    end
end

function grant.init()
    local ok, err = ddl.set_schema(grant.schema)
    if err then
        error(err)
    end

    box.schema.func.create('delete_grants_by_project_id')
    box.schema.user.grant('go', 'execute', 'function', 'delete_grants_by_project_id', {if_not_exists=true})
end

return grant