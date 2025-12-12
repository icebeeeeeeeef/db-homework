local key = KEYS[1]
--用户输入的验证码
local expectedCode = ARGV[1]
--redis中存储的验证码
local cntKey=key .. ":cnt"
--还需要验证次数是否用完
local cnt= tonumber(redis.call("GET",cntKey))

local code=redis.call("GET",key)

if cnt==nil then
    return -3
end

if cnt<=0 then
    return -1
elseif expectedCode==code then
    redis.call("del",key,cntKey)
    return 0
else
    --还没超过次数，但是验证码不正确
    redis.call("decr",cntKey)
    return -2
end

