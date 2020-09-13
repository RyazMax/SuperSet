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

return grant