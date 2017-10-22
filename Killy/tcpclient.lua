json = require "json"

TCP_CONN = nil
TCP_DATA = ""

TCP_CLIENT = {

  OnConnected = function (TCPConn)
    -- The specified link has succeeded in connecting to the remote server.
    -- Only called if the link is being connected as a client (using cNetwork:Connect() )
    -- Not used for incoming server links
    -- All returned values are ignored
    LOG("tcp client connected")
    TCP_CONN = TCPConn

    -- list containers
    LOG("listing containers...")
    -- SendTCPMessage("info",{"containers"},0)
  end,

  OnError = function (TCPConn, ErrorCode, ErrorMsg)
    -- The specified error has occurred on the link
    -- No other callback will be called for this link from now on
    -- For a client link being connected, this reports a connection error (destination unreachable etc.)
    -- It is an Undefined Behavior to send data to a_TCPLink in or after this callback
    -- All returned values are ignored
    LOG("tcp client OnError: " .. ErrorCode .. ": " .. ErrorMsg)

    -- retry to establish connection
    LOG("retry cNetwork:Connect")
    cNetwork:Connect("10.1.4.12",25566,TCP_CLIENT)
  end,

  OnReceivedData = function (TCPConn, Data)
    -- Data has been received on the link
    -- Will get called whenever there's new data on the link
    -- a_Data contains the raw received data, as a string
    -- All returned values are ignored
    -- LOG("TCP_CLIENT OnReceivedData")

    TCP_DATA = TCP_DATA .. Data
    local shiftLen = 0

    for message in string.gmatch(TCP_DATA, '([^\n]+\n)') do
      shiftLen = shiftLen + string.len(message)
      -- remove \n at the end
      message = string.sub(message,1,string.len(message)-1)
      ParseTCPMessage(message)
    end

    TCP_DATA = string.sub(TCP_DATA,shiftLen+1)

  end,

  OnRemoteClosed = function (TCPConn)
    -- The remote peer has closed the link
    -- The link is already closed, any data sent to it now will be lost
    -- No other callback will be called for this link from now on
    -- All returned values are ignored
    LOG("tcp client OnRemoteClosed")

    -- retry to establish connection
    LOG("retry cNetwork:Connect")
    cNetwork:Connect("10.1.4.12",25566,TCP_CLIENT)
  end,
}

-- SendTCPMessage sends a message over global
-- tcp connection TCP_CONN. args and id are optional
-- id stands for the request id.
function SendTCPMessage(cmd, args, data,id)
  if TCP_CONN == nil
  then
    LOG("can't send TCP message, TCP_CLIENT not connected")
    return
  end
  local v = {cmd=cmd,args={args},data=data,id=id}
  local msg = json.stringify(v) .. "\n"
  LOG(msg)
  TCP_CONN:Send(msg)
end

-- ParseTCPMessage parses a message received from
-- global tcp connection TCP_CONN
function ParseTCPMessage(message)
  local m = json.parse(message)
  -- deal with table events
  if m.cmd == "event" and table.getn(m.args) > 0 and m.args[1] == "table"
  then
    handleTableEvent(m.data)
  -- deal with monitor events
  elseif m.cmd == "monitor" and table.getn(m.args) > 0 and m.args[1] == "all"
  then
    handleMonitorEvent(m.data)
  elseif m.cmd == "event" and table.getn(m.args) > 0 and m.args[1] == "error"
  then
    localPlayer:SendMessage(cCompositeChat()
		:AddTextPart(m.data,"@c"))
  elseif m.cmd == "event" and table.getn(m.args) > 0 and m.args[1] == "result"
  then
    handleQueryEvent(m.data)
  end
end

function handleQueryEvent(event)
  LOG("handleQueryEvent")
  LOG(tostring(table.getn(event)))
  for i=1, table.getn(event.data)
  do
    LOG(tostring(event.data[i]))
    updateTableRecordContainer(event.data[i][1], "query", event.data[i])
    localPlayer:SendMessage(cCompositeChat()
		:AddTextPart(table.concat(event.data[i], " ")))
  end
  LOG("handleQueryEvent End")
end

function handleTableEvent(event)
  LOG("handleTableEvent")
  updateTableRecordContainer(event[1].name, "column", event[1].columns)
  for i=1, table.getn(event[1].data)
  do
    updateTableRecordContainer(event[1].data[i][1], tostring(i), event[1].data[i])
  end
  TABLE_SIGNAL_OFFSET = 10
  LOG("handleTableEvent End")
end

function handleMonitorEvent(event)
  LOG("handleMonitorEvent")
  for i=1, table.getn(event.TidbAvailHosts)
  do
    LOG("TidbAvailHosts")
    updateActiveInstanceContainer(event.TidbAvailHosts[i], "TiDB Instance", event.TidbAvailHosts[i], true)
  end
  for i=1, table.getn(event.TidbUnavailHosts)
  do
    LOG("TidbUnavailHosts")
    updateActiveInstanceContainer(event.TidbUnavailHosts[i], "TiDB Instance", event.TidbUnavailHosts[i], false)
  end
  for i=1, table.getn(event.TikvAvailHosts)
  do
    updateActiveInstanceContainer(event.TikvAvailHosts[i], "TiKV Instance", event.TikvAvailHosts[i], true)
  end
  for i=1, table.getn(event.TikvUnavailHosts)
  do
    updateActiveInstanceContainer(event.TikvUnavailHosts[i], "TiKV Instance", event.TikvUnavailHosts[i], false)
  end
  for i=1, table.getn(event.PdAvailHosts)
  do
    updateActiveInstanceContainer(event.PdAvailHosts[i], "PD Instance", event.PdAvailHosts[i], true)
  end
  for i=1, table.getn(event.PdUnavailHosts)
  do
    updateActiveInstanceContainer(event.PdUnavailHosts[i], "PD Instance", event.PdUnavailHosts[i], false)
  end
  for i=1, table.getn(event.EveryTikvStatus)
  do
    updateActiveInstanceStats(event.EveryTikvStatus[i].address, event.EveryTikvStatus[i].capacity, event.EveryTikvStatus[i].available)
  end
  LOG("handleMonitorEvent End")
end
