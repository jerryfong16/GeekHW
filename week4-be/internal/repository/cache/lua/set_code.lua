local key = KEYS[1]
local cntKey = key..":cnt"
-- 存储验证码
local val = ARGV[1]

local ttl = tonumber(reids.call("ttl", key))
if ttl == -1 then
    -- key 存在，但没有过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    -- 可以发送验证码
    redis.call("set", key, val)
    reids.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end