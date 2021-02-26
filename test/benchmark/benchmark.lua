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
    roll = random(6)
    if roll == 1 then
      wrk.body = '{"id":"' .. uuid() .. '",' .. '"name": "Pat Smith","email": "email' .. wrk.thread:get('id') .. '+' .. counter .. '"}'
      counter = counter + 1
      wrk.path = '/users'
      wrk.method = 'POST'
    elseif roll == 2 then
      wrk.path = '/tasks'
      wrk.headers["uid"] = "1490c27c-9698-49a4-bc5d-2f84b12ac95b"
      wrk.method = 'GET'
    elseif roll == 3 then
      wrk.body = '{"name":"' .. uuid() .. '","taskType": "basic"}'
      wrk.path = '/tasks'
      wrk.headers["uid"] = "1490c27c-9698-49a4-bc5d-2f84b12ac95b"
      wrk.method = 'POST'
    else
      wrk.path = '/hc'
      wrk.method = 'GET'
    end
    return wrk.format(nil, path)
end
