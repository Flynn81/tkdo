math.randomseed(os.time())
math.random(); math.random(); math.random()

random = math.random
function uuid()
    local template ='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function (c)
        local v = (c == 'x') and random(0, 0xf) or random(8, 0xb)
        return string.format('%x', v)
    end)
end

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json; charset=utf-8"

counter = 0
threadCounter = 0

setup = function(thread)
  thread:set('id',threadCounter)
  threadCounter = threadCounter + 1
end

request = function()
    wrk.body = '{"id":"' .. uuid() .. '",' .. '"name": "Pat Smith","email": "email' .. wrk.thread:get('id') .. '+' .. counter .. '"}'
    counter = counter + 1
    return wrk.format(nil, path)
end
