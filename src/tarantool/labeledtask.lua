local ddl = require('ddl')

local labeledtask = {}

labeledtask.schema = {
    spaces = {
        labeledtasks = {
            engine = 'memtx',
            is_local = false,
            temporary = false,
            format = {
                {name = 'ID', is_nullable = false, type = 'unsigned'},
                {name = 'ProjectID', is_nullable = false, type = 'unsigned'},
                {name = 'OriginID', is_nullable = false, type = 'unsigned'},
                {name = 'AnswerJSON', is_nullable = false, type = 'string'},
                {name = 'Timestamp', is_nullable = false, type = 'unsigned'},
            },
            indexes = {{
                name = 'ID',
                type = 'HASH',
                unique = true,
                parts = {
                    {path = 'ID', is_nullable = false, type = 'unsigned'},
                }
            },{
                name = 'ProjectID',
                type = 'TREE',
                unique = false,
                parts = {
                    {path = 'ProjectID', is_nullable = false, type = 'unsigned'},
                }
            },{
                name = 'OriginID',
                type = 'TREE',
                unique = false,
                parts = {
                    {path = 'ProjectID', is_nullable = false, type = 'unsigned'},
                    {path = 'OriginID', is_nullable = false, type = 'unsigned'},
                }
            },{
                name = 'Timestamp',
                type = 'TREE',
                unique = false,
                parts = {
                    {path = 'ProjectID', is_nullable = false, type = 'unsigned'},
                    {path = 'Timestamp', is_nullable = false, type = 'unsigned'}
                }
            }
        },
        },
    },
}

function task_greater_ts(pid, ts)
    local res = {}
    local count = 1
    for _,t in box.space.labeledtasks.index.Timestamp:pairs({pid, ts}, {iterator = 'GE'}) do
        if t['ProjectID'] ~= pid then
            break
        elseif ts < t['Timestamp'] then
            res[count] = t
            count = count + 1
        end
    end
    return res
end

function labeledtask.init()
    local ok, err = ddl.set_schema(labeledtask.schema)
    if err ~= nil then
        error(err)
    end

    box.schema.func.create('task_greater_ts', {if_not_exists=true})
    box.schema.user.grant('go', 'execute', 'function', 'task_greater_ts', {if_not_exists=true})
end

return labeledtask