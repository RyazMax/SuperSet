local schema = {}

schema.schema = {
    spaces = {
        project_schemas = {
            engine = 'memtx',
            is_local = false,
            temporary = false,
            format = {
                {name = 'ProjectID', is_nullable = false, type = 'unsigned'},
                {name = 'InputType', is_nullable = false, type = 'string'},
                {name = 'InputSchema', is_nullable = false, type = 'string'},
                {name = 'OutputType', is_nullable = false, type = 'string'},
                {name = 'OutputSchema', is_nullable = false, type = 'string'},
            },
            indexes = {{
                name = 'ID',
                type = 'HASH',
                unique = true,
                parts = {
                    {path = 'ProjectID', is_nullable = false, type = 'unsigned'}
                }
            },
        },
        },
    },
}

return schema