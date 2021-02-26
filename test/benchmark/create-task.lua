math.randomseed(os.time())
math.random(); math.random(); math.random()

random = math.random
function uuid()
    local template ='00000000-0000-0000-0000-000000000000'
    return string.gsub(template, '[xy]', function (c)
        local v = (c == 'x') and random(0, 0xf) or random(8, 0xb)
        return string.format('%x', v)
    end)
end

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json; charset=utf-8"
wrk.headers["uid"] = "1490c27c-9698-49a4-bc5d-2f84b12ac95b"

counter = 0
threadCounter = 0

setup = function(thread)
  thread:set('id',threadCounter)
  threadCounter = threadCounter + 1
end

request = function()
    wrk.body = '{"name":"' .. uuid() .. '","taskType": "basic"}'
    counter = counter + 1
    return wrk.format(nil, path)
end
