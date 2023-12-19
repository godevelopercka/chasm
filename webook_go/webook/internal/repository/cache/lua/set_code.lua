-- 获取传入的 Redis key，这个 key 是用于存储验证码的
local key = KEYS[1]
-- 构建一个用于存储验证码验证次数的 key
-- phone_code：login:152xxxxxxx:cnt
local cntKey = key..":cnt"
-- 你的验证码 123456
-- 获取传入的验证码值
local val = ARGV[1]
-- 获取验证码的过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- 如果过期时间为 -1，说明这个 key 存在但没有设置过期时间
    -- 这可能是一个系统错误，例如同事手动设置了 key 但没有设置过期时间
    return -2
    -- 如果过期时间小于 540（600秒 - 60秒的缓冲时间），说明验证码被频繁尝试
    -- 这种情况可能是由于网络延迟、重试、恶意攻击等原因导致的
    -- 在这种情况下，我们重新设置验证码和过期时间，并重置验证次数为3
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val) -- 重新设置验证码
    redis.call("expire", key, 600) -- 设置新的过期时间为600秒
    redis.call("set", cntKey, 3) -- 重置验证次数为3
    redis.call("expire", cntKey, 600) -- 设置新的过期时间为600秒
    -- 完美，符合预期
    return 0 -- 返回0，表示操作成功
else
    -- 如果过期时间大于或等于540秒，说明验证码发送太频繁
    -- 这种情况下，我们不进行任何操作，并返回-1作为错误码
    return -1
end