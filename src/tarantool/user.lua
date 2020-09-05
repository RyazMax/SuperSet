local user = {}

user.schema = { 
    spaces = {
        users = {
            engine = 'memtx',
            is_local = false,
            temporary = false,
            format = {
                {name = 'ID', is_nullable = false, type = 'unsigned'},
                {name = 'Login', is_nullable = false, type = 'string'},
                {name = 'PasswordShadow', is_nullable = false, type = 'string'},
                {name = 'Type', is_nullable = false, type = 'unsigned'},
            },
            indexes = {{
                name = 'ID',
                type = 'HASH',
                unique = true,
                parts = {
                    {path = 'ID', is_nullable = false, type = 'unsigned'}
                }
            }, {
                name = 'login',
                type = 'TREE',
                unique = true,
                parts = {
                    {path = 'Login', is_nullable = false, type = 'string'}
                }
            }},
        },
    },
}

return user