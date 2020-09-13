local session = {}

session.schema = {
    spaces = {
        sessions = {
            engine = 'memtx',
            is_local = false,
            temporary = false,
            format = {
                {name = 'ID', is_nullable = false, type = 'string'},
                {name = 'UserLogin', is_nullable = false, type = 'string'},
            },
            indexes = {{
                name = 'ID',
                type = 'HASH',
                unique = true,
                parts = {
                    {path = 'ID', is_nullable = false, type = 'string'}
                }
            }, {
                name = 'userlogin',
                type = 'TREE',
                unique = false,
                parts = {
                    {path = 'UserLogin', is_nullable = false, type = 'string'}
                }
            }},
        },
    },
}

return session