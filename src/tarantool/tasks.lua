local queue = require('queue')

local tasks = {}
local TTR = 60 * 15 -- 15 Minutes
local ITERATIONS_NUM = 10
local PER_PROJECT_TIMEOUT = 0.1

local function take_task_by_projects(...)
    local projs = {...}
    local task
    for i = 1,ITERATIONS_NUM do
        for _, proj in ipairs(projs[1]) do
            task = queue.tube[proj]:take(PER_PROJECT_TIMEOUT)
            if task ~= nil then
                return task
            end
        end
    end
end

function take_aggr_by_projects(...)
    local task = take_task_by_projects(...)
    if task == nil then
        return nil
    end

    return {task[1], task[3]} -- DataPart and ID part
end

function ack_task(proj, id)
    queue.tube[proj]:ack(id)
end

function create_tube(proj)
    queue.create_tube(proj, 'fifottl')
end

function drop_tube(proj)
    queue.tube[proj]:drop()
end

function insert_task(proj, task) 
    queue.tube[proj]:put(task, {ttr = TTR})
end

local func_names = {
    'insert_task',
    'create_tube',
    'drop_tube',
    'ack_task',
    'take_aggr_by_projects',
}

function tasks.init()
    for _, name in ipairs(func_names) do
        box.schema.func.create(name, {if_not_exists=true})
        box.schema.user.grant('go', 'execute', 'function', name, {if_not_exists=true})
    end
    box.schema.user.grant('go', 'super', nil, nil, {if_not_exists=true})
end

return tasks