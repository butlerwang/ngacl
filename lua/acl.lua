local http = require "resty.http"
local httpc = http.new()
local uri = ngx.var.request_uri
local host = ngx.var.host

local function isempty(s)
  return s == nil or s == ''
end

local failover = ngx.var.ccna_acl_failover

if failover == '' then
  failover = 200
end
failover = tonumber(failover)

httpc:set_timeout(20)
httpc:connect("127.0.0.1", 8104)

local h = ngx.req.get_headers()
local header = {}
for k, v in pairs(h) do
  if k:lower() ~= "content-length" and k:lower() ~= "connection" then
    header[k] = v:gsub("[^\x20-\x7E]","")
  end
end

header["host"] = host
header["ccna-acl-remote-addr"] = ngx.var.remote_addr
header["ccna-acl-request-method"] = ngx.var.request_method
header["ccna-acl-method"] = ngx.var.ccna_acl_method
header["ccna-acl-scheme"] = ngx.var.scheme

-- start hmac
if not isempty(ngx.var.ccna_acl_hmac_key1) then
  header["ccna-acl-hmac-key1"] = ngx.var.ccna_acl_hmac_key1
end
if not isempty(ngx.var.ccna_acl_hmac_key2) then
  header["ccna-acl-Hmac-key2"] = ngx.var.ccna_acl_hmac_key2
end
-- end hmac
-- start tkmd5
if not isempty(ngx.var.ccna_acl_tkmd5_secret) then
  header["ccna-acl-tkmd5-secret"] = ngx.var.ccna_acl_tkmd5_secret
end
-- end tkmd5

-- start s3
if not isempty(ngx.var.ccna_s3_accesskey) then
  header["ccna-s3-accesskey"] = ngx.var.ccna_s3_accesskey
end
if not isempty(ngx.var.ccna_s3_secretkey) then
  header["ccna-s3-secretkey"] = ngx.var.ccna_s3_secretkey
end
if not isempty(ngx.var.ccna_s3_region) then
  header["ccna-s3-region"] = ngx.var.ccna_s3_region
end
if not isempty(ngx.var.ccna_s3_bucket) then
  header["ccna-s3-bucket"] = ngx.var.ccna_s3_bucket
end
-- end s3


local res, err = httpc:request{
    path = uri,
    headers = header,
}

if err then
  ngx.log(ngx.ERR, err)
  if failover >= 400 then
    ngx.exit(failover)
  end
  return
end
if not res then
  ngx.log(ngx.ERR, "failed to request")
  httpc:set_keepalive()
  if failover >= 400 then
    ngx.exit(failover)
  end
  return
end

ngx.header["X-Cc-Auth-Id"] = res.headers["X-Cc-Auth-Id"]

-- start s3
if not isempty(res.headers["Authorization"]) then
  ngx.req.set_header("Authorization",res.headers["Authorization"])
end
if not isempty(res.headers["X-Amz-Content-Sha256"]) then
  ngx.req.set_header("X-Amz-Content-Sha256",res.headers["X-Amz-Content-Sha256"])
end
if not isempty(res.headers["X-Amz-Date"]) then
  ngx.req.set_header("X-Amz-Date",res.headers["X-Amz-Date"])
end
-- end s3

if res.status >= 400 then
  httpc:set_keepalive()
  ngx.exit(res.status)
end

httpc:set_keepalive()
