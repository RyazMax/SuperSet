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

return project