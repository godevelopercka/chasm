-- 从 Redis 中获取传入的 key
local key = KEYS[1]
-- 构建一个用于存储验证码验证次数的 key
local cntKey = key .. ":cnt"
-- 用户输入的 code
local expectedCode = ARGV[1]
-- 从 Redis 中获取与 key 对应的验证码
local code = redis.call("get", key)
-- 从 Redis 中获取与 cntKey 对应的验证次数
local cnt = tonumber(redis.call("get", cntKey))
-- 如果验证次数小于或等于0
if cnt <= 0 then
-- 说明用户一直输错验证码，或者验证码已经被使用过了
-- 可能有人恶意尝试，返回错误码 -1
    return -1
end
-- 如果用户输入的验证码与 Redis 中存储的验证码相同
if code == expectedCode then
    -- 说明验证码输入正确
    -- 将验证次数设置为 -1，表示该验证码已经使用过
    redis.call("set", cntKey, -1)
    return 0
else
    -- 如果用户输入的验证码与 Redis 中存储的验证码不同
    -- 验证次数减1
    redis.call("decr", cntKey)
    return -2
end